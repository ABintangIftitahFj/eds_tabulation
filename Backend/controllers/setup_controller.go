package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star_fj/eds-backend/models"
)

// 1. Buat Turnamen Baru
func CreateTournament(c *gin.Context) {
	var input models.Tournament
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tournament: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": input})
}

// 1b. Update Turnamen (Edit Status, dll)
func UpdateTournament(c *gin.Context) {
	id := c.Param("id")
	var tournament models.Tournament
	if err := models.DB.First(&tournament, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	var input models.Tournament
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Update fields
	tournament.Status = input.Status
	if input.Name != "" {
		tournament.Name = input.Name
	}
	if input.Location != "" {
		tournament.Location = input.Location
	}
	models.DB.Save(&tournament)
	c.JSON(http.StatusOK, gin.H{"data": tournament})
}

func DeleteTournament(c *gin.Context) {
	id := c.Param("id")
	var tournament models.Tournament
	if err := models.DB.First(&tournament, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	if err := models.DB.Delete(&tournament).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tournament deleted successfully"})
}

// 2. Daftarkan Tim ke Turnamen
func CreateTeam(c *gin.Context) {
	var input models.Team
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": input})
}

func DeleteTeam(c *gin.Context) {
	id := c.Param("id")
	var team models.Team
	if err := models.DB.First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}
	if err := models.DB.Delete(&team).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
}

// 3. Buat Ronde (Round 1, 2, dll)
func CreateRound(c *gin.Context) {
	var input models.Round
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create round: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": input})
}

func DeleteRound(c *gin.Context) {
	id := c.Param("id")
	var round models.Round
	if err := models.DB.First(&round, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Round not found"})
		return
	}
	if err := models.DB.Delete(&round).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Round deleted successfully"})
}

// 5. LIHAT DAFTAR TURNAMEN (Baru)
func GetTournaments(c *gin.Context) {
	var tournaments []models.Tournament
	models.DB.Order("created_at desc").Find(&tournaments)
	c.JSON(http.StatusOK, gin.H{"data": tournaments})
}

// 6. LIHAT DAFTAR TIM (Baru)
func GetTeams(c *gin.Context) {
	tournamentID := c.Query("tournament_id")
	var teams []models.Team
	query := models.DB.Preload("Speakers").Order("created_at desc")
	if tournamentID != "" {
		query = query.Where("tournament_id = ?", tournamentID)
	}
	if err := query.Find(&teams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if teams == nil {
		teams = []models.Team{}
	}
	c.JSON(http.StatusOK, gin.H{"data": teams})
}

// 7. LIHAT DAFTAR RONDE (Berdasarkan Turnamen ID)
func GetRounds(c *gin.Context) {
	tournamentID := c.Query("tournament_id")
	var rounds []models.Round
	query := models.DB.Order("created_at asc")
	if tournamentID != "" {
		query = query.Where("tournament_id = ?", tournamentID)
	}
	if err := query.Find(&rounds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rounds == nil {
		rounds = []models.Round{}
	}
	c.JSON(http.StatusOK, gin.H{"data": rounds})
}

// 9. Buat Match (Pairing: Tim A vs Tim B)
func CreateMatch(c *gin.Context) {
	var input struct {
		RoundID     uint   `json:"round_id"`
		GovTeamID   uint   `json:"gov_team_id"`
		OppTeamID   uint   `json:"opp_team_id"`
		Room        string `json:"room"`
		Adjudicator string `json:"adjudicator"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Debug log
	fmt.Printf("CreateMatch received: %+v\n", input)
	match := models.Match{
		RoundID:     input.RoundID,
		GovTeamID:   &input.GovTeamID,
		OppTeamID:   &input.OppTeamID,
		Room:        input.Room,
		Adjudicator: input.Adjudicator,
		IsCompleted: false,
	}
	if err := models.DB.Create(&match).Error; err != nil {
		println("CreateMatch DB error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": match})
}

// 10. LIHAT DAFTAR MATCH (Pairing)
func GetMatches(c *gin.Context) {
	roundID := c.Query("round_id") // Filter per ronde
	var matches []models.Match
	query := models.DB.Preload("GovTeam").Preload("OppTeam").Preload("Round").Order("room asc")
	if roundID != "" {
		query = query.Where("round_id = ?", roundID)
	}
	if err := query.Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if matches == nil {
		matches = []models.Match{}
	}
	c.JSON(http.StatusOK, gin.H{"data": matches})
}

// 11. UPDATE MATCH RESULT (Set Winner)
func UpdateMatchResult(c *gin.Context) {
	matchID := c.Param("id")
	var input struct {
		WinnerID    uint `json:"winner_id"`
		IsCompleted bool `json:"is_completed"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var match models.Match
	if err := models.DB.First(&match, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	// Update match result
	match.WinnerID = &input.WinnerID
	match.IsCompleted = input.IsCompleted

	if err := models.DB.Save(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": match})
}

// 12. ADJUDICATOR MANAGEMENT
func CreateAdjudicator(c *gin.Context) {
	var input models.Adjudicator
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": input})
}

func GetAdjudicators(c *gin.Context) {
	tournamentID := c.Query("tournament_id")
	var adjudicators []models.Adjudicator
	query := models.DB.Order("name asc")
	if tournamentID != "" {
		query = query.Where("tournament_id = ?", tournamentID)
	}
	if err := query.Find(&adjudicators).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if adjudicators == nil {
		adjudicators = []models.Adjudicator{}
	}
	c.JSON(http.StatusOK, gin.H{"data": adjudicators})
}

// 13. ROOM MANAGEMENT
func CreateRoom(c *gin.Context) {
	var input models.Room
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": input})
}

func GetRooms(c *gin.Context) {
	tournamentID := c.Query("tournament_id")
	var rooms []models.Room
	query := models.DB.Order("name asc")
	if tournamentID != "" {
		query = query.Where("tournament_id = ?", tournamentID)
	}
	if err := query.Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rooms == nil {
		rooms = []models.Room{}
	}
	c.JSON(http.StatusOK, gin.H{"data": rooms})
}

// DELETE ENDPOINTS
func DeleteAdjudicator(c *gin.Context) {
	id := c.Param("id")
	var adjudicator models.Adjudicator
	if err := models.DB.First(&adjudicator, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Adjudicator not found"})
		return
	}
	if err := models.DB.Delete(&adjudicator).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Adjudicator deleted successfully"})
}

func DeleteRoom(c *gin.Context) {
	id := c.Param("id")
	var room models.Room
	if err := models.DB.First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}
	if err := models.DB.Delete(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})
}

func DeleteMatch(c *gin.Context) {
	id := c.Param("id")
	var match models.Match
	if err := models.DB.First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}
	if err := models.DB.Delete(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Match deleted successfully"})
}
