package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star_fj/eds-backend/models"
)

type BallotInput struct {
	MatchID     uint            `json:"match_id"`
	Adjudicator string          `json:"adjudicator"`
	Scores      []models.Ballot `json:"scores"`
}

func SubmitBallot(c *gin.Context) {
	var input BallotInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mulai Transaksi Database (Biar Aman)
	tx := models.DB.Begin()

	// 0. Ambil data match terlebih dahulu untuk mendapat teamID
	var match models.Match
	if err := tx.Preload("GovTeam").Preload("OppTeam").First(&match, input.MatchID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Match tidak ditemukan"})
		return
	}

	// 1. Simpan Skor Individu
	var totalGov float64 = 0
	var totalOpp float64 = 0

	for _, ballot := range input.Scores {
		ballot.MatchID = input.MatchID

		// Validasi data pembicara (Jaga-jaga)
		if ballot.Speaker.Name == "" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nama speaker tidak boleh kosong"})
			return
		}

		// Tentukan TeamID berdasarkan TeamRole dan MatchID
		var teamID uint
		if ballot.TeamRole == "gov" {
			if match.GovTeamID == nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Match tidak memiliki Government Team"})
				return
			}
			teamID = *match.GovTeamID
		} else if ballot.TeamRole == "opp" {
			if match.OppTeamID == nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Match tidak memiliki Opposition Team"})
				return
			}
			teamID = *match.OppTeamID
		} else {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "TeamRole harus 'gov' atau 'opp'"})
			return
		}

		// Debug log
		fmt.Printf("Debug: TeamRole=%s, TeamID=%d, MatchID=%d\n", ballot.TeamRole, teamID, ballot.MatchID)

		// Validasi TeamID tidak boleh 0
		if teamID == 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "TeamID tidak valid (0)"})
			return
		}

		// Cari atau buat speaker baru jika belum ada
		var speaker models.Speaker
		if ballot.SpeakerID != 0 {
			// Jika SpeakerID sudah ada, gunakan yang ada
			if err := tx.First(&speaker, ballot.SpeakerID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusNotFound, gin.H{"error": "Speaker tidak ditemukan"})
				return
			}
		} else if ballot.Speaker.Name != "" {
			// Cari speaker berdasarkan nama dan tim
			err := tx.Where("name = ? AND team_id = ?", ballot.Speaker.Name, teamID).First(&speaker).Error
			if err != nil {
				// Jika tidak ditemukan, buat speaker baru
				speaker = models.Speaker{
					Name:   ballot.Speaker.Name,
					TeamID: teamID,
				}
				if err := tx.Create(&speaker).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat speaker: " + err.Error()})
					return
				}
			}
			ballot.SpeakerID = speaker.ID
		} else {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "SpeakerID atau Speaker.Name harus diisi"})
			return
		}

		// Simpan ballot dengan SpeakerID yang valid
		ballotToSave := models.Ballot{
			MatchID:   ballot.MatchID,
			SpeakerID: ballot.SpeakerID,
			Score:     ballot.Score,
			Position:  ballot.Position,
			IsReply:   ballot.IsReply,
			TeamRole:  ballot.TeamRole,
		}

		if err := tx.Create(&ballotToSave).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan skor: " + err.Error()})
			return
		}

		// Hitung Total Skor
		if ballot.TeamRole == "gov" {
			totalGov += ballot.Score
		} else if ballot.TeamRole == "opp" {
			totalOpp += ballot.Score
		}
	}

	// 2. Tentukan Pemenang (Logic Asian Parliamentary)
	// Match sudah diambil di atas, tidak perlu diambil lagi

	// Update Status Match
	match.IsCompleted = true
	// match.Adjudicator = input.Adjudicator // Removed because Adjudicator is now a relation

	// Logika Penentuan Pemenang
	// Kalau Gov > Opp -> Gov Menang (WinnerID = GovTeamID)
	// Kalau Opp > Gov -> Opp Menang
	// Kalau Seri -> (Di debat jarang seri, biasanya juri dipaksa milih margin tipis)
	if totalGov > totalOpp {
		match.WinnerID = match.GovTeamID
	} else {
		match.WinnerID = match.OppTeamID
	}

	if err := tx.Save(&match).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update hasil match"})
		return
	}

	// 3. Update Klasemen Tim (Standings)
	// Tambahkan VP (Victory Point) dan Speaker Score ke masing-masing tim

	// Update Gov Team
	match.GovTeam.TotalSpeaker += totalGov
	if match.WinnerID == match.GovTeamID {
		match.GovTeam.TotalVP += 1 // Menang dapat 1 poin
		match.GovTeam.Wins += 1
	} else {
		match.GovTeam.Losses += 1
	}
	if err := tx.Save(&match.GovTeam).Error; err != nil {
		tx.Rollback()
		return
	}

	// Update Opp Team
	match.OppTeam.TotalSpeaker += totalOpp
	if match.WinnerID == match.OppTeamID {
		match.OppTeam.TotalVP += 1
		match.OppTeam.Wins += 1
	} else {
		match.OppTeam.Losses += 1
	}
	if err := tx.Save(&match.OppTeam).Error; err != nil {
		tx.Rollback()
		return
	}

	// Selesai!
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message":   "Skor disimpan & Pemenang ditentukan!",
		"winner_id": match.WinnerID,
		"total_gov": totalGov,
		"total_opp": totalOpp,
	})
}

func GetBallots(c *gin.Context) {
	matchID := c.Query("match_id")
	roundID := c.Query("round_id")

	var ballots []models.Ballot

	if matchID != "" {
		// Query by specific match
		if err := models.DB.Preload("Speaker").Where("match_id = ?", matchID).Find(&ballots).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if roundID != "" {
		// Query by round - get all match IDs first
		var matchIDs []uint
		if err := models.DB.Table("matches").Where("round_id = ?", roundID).Pluck("id", &matchIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get matches: " + err.Error()})
			return
		}

		if len(matchIDs) == 0 {
			// No matches found, return empty array
			c.JSON(http.StatusOK, gin.H{"data": []models.Ballot{}})
			return
		}

		if err := models.DB.Preload("Speaker").Where("match_id IN ?", matchIDs).Find(&ballots).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "match_id or round_id is required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ballots})
}
