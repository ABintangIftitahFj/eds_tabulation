package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdjudicatorFeedback struct {
	ID            uint      `json:"ID" gorm:"primaryKey"`
	MatchID       uint      `json:"match_id"`
	TournamentID  uint      `json:"tournament_id"`
	AdjudicatorID uint      `json:"adjudicator_id"`
	Rating        int       `json:"rating"`  // 1-5 stars
	Comment       *string   `json:"comment"` // Optional comment
	CreatedAt     time.Time `json:"created_at"`
}

// GetAdjudicatorFeedback - Get feedbacks with optional filters
func GetAdjudicatorFeedback(c *gin.Context, db *gorm.DB) {
	var feedbacks []AdjudicatorFeedback
	query := db.Model(&AdjudicatorFeedback{})

	// Filter by adjudicator_id if provided
	if adjID := c.Query("adjudicator_id"); adjID != "" {
		query = query.Where("adjudicator_id = ?", adjID)
	}

	// Filter by tournament_id if provided
	if tourID := c.Query("tournament_id"); tourID != "" {
		query = query.Where("tournament_id = ?", tourID)
	}

	// Filter by match_id if provided
	if matchID := c.Query("match_id"); matchID != "" {
		query = query.Where("match_id = ?", matchID)
	}

	if err := query.Order("created_at DESC").Find(&feedbacks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch feedbacks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": feedbacks})
}

// CreateAdjudicatorFeedback - Create new feedback
func CreateAdjudicatorFeedback(c *gin.Context, db *gorm.DB) {
	var feedback AdjudicatorFeedback

	if err := c.ShouldBindJSON(&feedback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if feedback already exists for this match
	var existingCount int64
	db.Model(&AdjudicatorFeedback{}).
		Where("match_id = ?", feedback.MatchID).
		Count(&existingCount)

	if existingCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah memberikan feedback untuk pertandingan ini"})
		return
	}

	// Validate rating
	if feedback.Rating < 1 || feedback.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
		return
	}

	// Require comment for low ratings (1-2)
	if feedback.Rating <= 2 && (feedback.Comment == nil || *feedback.Comment == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment is required for low ratings"})
		return
	}

	feedback.CreatedAt = time.Now()

	if err := db.Create(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create feedback"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": feedback})
}

// CheckFeedbackExists - Check if feedback already exists for a match
func CheckFeedbackExists(c *gin.Context, db *gorm.DB) {
	matchID := c.Query("match_id")

	if matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "match_id is required"})
		return
	}

	var count int64
	db.Model(&AdjudicatorFeedback{}).
		Where("match_id = ?", matchID).
		Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"exists": count > 0,
		"message": func() string {
			if count > 0 {
				return "Anda sudah memberikan feedback untuk pertandingan ini"
			}
			return "Belum ada feedback"
		}(),
	})
}

// GetFeedbackStats - Get statistics for an adjudicator
func GetFeedbackStats(c *gin.Context, db *gorm.DB) {
	adjID := c.Param("adjudicator_id")

	// Count total feedbacks
	var totalCount int64
	db.Model(&AdjudicatorFeedback{}).Where("adjudicator_id = ?", adjID).Count(&totalCount)

	// Calculate average rating
	var avgRating float64
	db.Model(&AdjudicatorFeedback{}).
		Where("adjudicator_id = ?", adjID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avgRating)

	// Rating distribution
	type RatingCount struct {
		Rating int   `json:"rating"`
		Count  int64 `json:"count"`
	}
	var distribution []RatingCount
	db.Model(&AdjudicatorFeedback{}).
		Select("rating, COUNT(*) as count").
		Where("adjudicator_id = ?", adjID).
		Group("rating").
		Order("rating DESC").
		Scan(&distribution)

	c.JSON(http.StatusOK, gin.H{
		"total_feedbacks":     totalCount,
		"average_rating":      avgRating,
		"rating_distribution": distribution,
	})
}

// DeleteAdjudicatorFeedback - Delete a feedback (admin only ideally)
func DeleteAdjudicatorFeedback(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := db.Delete(&AdjudicatorFeedback{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete feedback"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feedback deleted successfully"})
}
