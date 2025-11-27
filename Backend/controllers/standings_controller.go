package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star_fj/eds-backend/models"
)

// GET /api/standings/teams?tournament_id=1
func GetStandings(c *gin.Context) {
	tournamentID := c.Query("tournament_id")
	var teams []models.Team

	if tournamentID != "" {
		// Get all teams that participated in matches for this tournament
		// This includes teams that might not be properly registered in teams table
		var govTeamIDs []uint
		var oppTeamIDs []uint

		// Get government team IDs from matches in rounds of this tournament
		models.DB.Table("matches").
			Joins("JOIN rounds ON matches.round_id = rounds.id").
			Where("rounds.tournament_id = ?", tournamentID).
			Pluck("DISTINCT gov_team_id", &govTeamIDs)

		// Get opposition team IDs from matches in rounds of this tournament
		models.DB.Table("matches").
			Joins("JOIN rounds ON matches.round_id = rounds.id").
			Where("rounds.tournament_id = ?", tournamentID).
			Pluck("DISTINCT opp_team_id", &oppTeamIDs)

		// Combine and deduplicate team IDs
		teamIDMap := make(map[uint]bool)
		for _, id := range govTeamIDs {
			teamIDMap[id] = true
		}
		for _, id := range oppTeamIDs {
			teamIDMap[id] = true
		}

		var participatingTeamIDs []uint
		for id := range teamIDMap {
			participatingTeamIDs = append(participatingTeamIDs, id)
		}

		if len(participatingTeamIDs) > 0 {
			// Get teams that actually participated
			models.DB.Where("id IN ?", participatingTeamIDs).
				Order("total_vp desc").Order("total_speaker desc").
				Find(&teams)
		} else {
			// Fallback to original method if no matches found
			models.DB.Where("tournament_id = ?", tournamentID).
				Order("total_vp desc").Order("total_speaker desc").
				Find(&teams)
		}
	} else {
		models.DB.Order("total_vp desc").Order("total_speaker desc").Find(&teams)
	}

	// Update Ranking Angka (1, 2, 3...) secara manual sebelum dikirim
	for i := range teams {
		teams[i].Rank = i + 1
	}

	// Handle empty results
	if teams == nil {
		teams = []models.Team{}
	}

	c.JSON(http.StatusOK, gin.H{"data": teams})
}

// GET /api/standings/speakers?tournament_id=1
func GetSpeakerStandings(c *gin.Context) {
	tournamentID := c.Query("tournament_id")
	var speakers []models.Speaker

	// Join with Team to filter by tournament_id
	query := models.DB.Joins("JOIN teams ON teams.id = speakers.team_id").
		Select("speakers.*, teams.name as team_name, teams.institution as institution").
		Order("speakers.total_score desc")

	if tournamentID != "" {
		query = query.Where("teams.tournament_id = ?", tournamentID)
	}

	if err := query.Find(&speakers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Assign ranks
	for i := range speakers {
		speakers[i].SpeakerRank = i + 1
	}

	if speakers == nil {
		speakers = []models.Speaker{}
	}

	c.JSON(http.StatusOK, gin.H{"data": speakers})
}

// GET /api/institutions?tournament_id=1
func GetParticipatingInstitutions(c *gin.Context) {
	tournamentID := c.Query("tournament_id")

	type InstitutionData struct {
		Institution string  `json:"institution"`
		TeamCount   int     `json:"team_count"`
		TotalPoints int     `json:"total_points"`
		AvgPoints   float64 `json:"avg_points"`
	}

	var institutions []InstitutionData

	if tournamentID != "" {
		// Get institutions from teams that participated in matches
		models.DB.Table("teams").
			Select("teams.institution, COUNT(DISTINCT teams.id) as team_count, SUM(teams.total_vp) as total_points, AVG(teams.total_vp) as avg_points").
			Joins("JOIN matches ON teams.id = matches.gov_team_id OR teams.id = matches.opp_team_id").
			Joins("JOIN rounds ON matches.round_id = rounds.id").
			Where("rounds.tournament_id = ?", tournamentID).
			Group("teams.institution").
			Order("total_points desc").
			Scan(&institutions)
	} else {
		// Get all institutions from teams table
		models.DB.Table("teams").
			Select("teams.institution, COUNT(teams.id) as team_count, SUM(teams.total_vp) as total_points, AVG(teams.total_vp) as avg_points").
			Group("teams.institution").
			Order("total_points desc").
			Scan(&institutions)
	}

	if institutions == nil {
		institutions = []InstitutionData{}
	}

	c.JSON(http.StatusOK, gin.H{"data": institutions})
}

// POST /api/standings/recalculate?tournament_id=1
// Menghitung ulang seluruh standings dari ballot yang ada
func RecalculateStandings(c *gin.Context) {
	tournamentID := c.Query("tournament_id")
	if tournamentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tournament_id diperlukan"})
		return
	}

	tx := models.DB.Begin()

	// 1. Reset semua team stats untuk tournament ini
	if err := tx.Model(&models.Team{}).
		Where("tournament_id = ?", tournamentID).
		Updates(map[string]interface{}{
			"total_vp":      0,
			"total_speaker": 0,
			"wins":          0,
			"losses":        0,
		}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal reset team stats"})
		return
	}

	// 2. Reset semua speaker scores untuk tournament ini
	if err := tx.Exec(`
		UPDATE speakers SET total_score = 0 
		WHERE team_id IN (SELECT id FROM teams WHERE tournament_id = ?)
	`, tournamentID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal reset speaker scores"})
		return
	}

	// 3. Ambil semua match yang sudah completed di tournament ini
	var completedMatches []models.Match
	if err := tx.Preload("GovTeam").Preload("OppTeam").
		Joins("JOIN rounds ON matches.round_id = rounds.id").
		Where("rounds.tournament_id = ? AND matches.is_completed = ?", tournamentID, true).
		Find(&completedMatches).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil completed matches"})
		return
	}

	fmt.Printf("Debug Recalculate: Found %d completed matches\n", len(completedMatches))

	// 4. Untuk setiap match, hitung ulang stats dari ballot
	for _, match := range completedMatches {
		// Ambil ballots untuk match ini
		var ballots []models.Ballot
		tx.Where("match_id = ?", match.ID).Find(&ballots)

		var totalGov, totalOpp int
		speakerScores := make(map[uint]int)

		for _, ballot := range ballots {
			if ballot.TeamRole == "gov" {
				totalGov += ballot.Score
			} else if ballot.TeamRole == "opp" {
				totalOpp += ballot.Score
			}
			if ballot.SpeakerID != 0 {
				speakerScores[ballot.SpeakerID] += ballot.Score
			}
		}

		// Update Gov Team
		if match.GovTeamID != nil {
			var govTeam models.Team
			if err := tx.First(&govTeam, *match.GovTeamID).Error; err == nil {
				govTeam.TotalSpeaker += totalGov
				if match.WinnerID != nil && *match.WinnerID == *match.GovTeamID {
					govTeam.TotalVP += 1
					govTeam.Wins += 1
				} else {
					govTeam.Losses += 1
				}
				tx.Save(&govTeam)
			}
		}

		// Update Opp Team
		if match.OppTeamID != nil {
			var oppTeam models.Team
			if err := tx.First(&oppTeam, *match.OppTeamID).Error; err == nil {
				oppTeam.TotalSpeaker += totalOpp
				if match.WinnerID != nil && *match.WinnerID == *match.OppTeamID {
					oppTeam.TotalVP += 1
					oppTeam.Wins += 1
				} else {
					oppTeam.Losses += 1
				}
				tx.Save(&oppTeam)
			}
		}

		// Update Speaker Scores
		for speakerID, score := range speakerScores {
			var speaker models.Speaker
			if err := tx.First(&speaker, speakerID).Error; err == nil {
				speaker.TotalScore += score
				tx.Save(&speaker)
			}
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message":           "Standings berhasil dihitung ulang",
		"matches_processed": len(completedMatches),
	})
}
