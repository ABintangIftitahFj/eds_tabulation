package controllers

import (
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
