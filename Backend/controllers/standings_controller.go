package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star_fj/eds-backend/models"
)

// GET /api/standings?tournament_id=1
func GetStandings(c *gin.Context) {
	tournamentID := c.Query("tournament_id")
	var teams []models.Team

	query := models.DB.Order("total_vp desc").Order("total_speaker desc")

	if tournamentID != "" {
		query = query.Where("tournament_id = ?", tournamentID)
	}

	// Ambil data tim yang sudah diurutkan berdasarkan VP tertinggi, lalu Skor tertinggi
	if err := query.Find(&teams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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
