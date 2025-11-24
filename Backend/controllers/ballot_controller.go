package controllers

import (
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

		// Simpan speaker baru jika belum ada (atau update)
		// Note: Di sistem real, idealnya speaker dipilih dari ID, tapi ini simplifikasi
		if err := tx.Create(&ballot).Error; err != nil {
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
	var match models.Match
	if err := tx.Preload("GovTeam").Preload("OppTeam").First(&match, input.MatchID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Match tidak ditemukan"})
		return
	}

	// Update Status Match
	match.IsCompleted = true
	match.Adjudicator = input.Adjudicator

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
