package controllers

import (
	"fmt"
	"net/http"
	"strconv"

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

// Publish/Unpublish Draw
func PublishDraw(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		IsDrawPublished bool `json:"is_draw_published"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var round models.Round
	if err := models.DB.First(&round, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Round not found"})
		return
	}

	round.IsDrawPublished = input.IsDrawPublished
	if err := models.DB.Save(&round).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	message := "Draw published to users"
	if !input.IsDrawPublished {
		message = "Draw hidden from users"
	}
	c.JSON(http.StatusOK, gin.H{"message": message, "data": round})
}

// Publish/Unpublish Motion
func PublishMotion(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		IsMotionPublished bool `json:"is_motion_published"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var round models.Round
	if err := models.DB.First(&round, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Round not found"})
		return
	}

	round.IsMotionPublished = input.IsMotionPublished
	if err := models.DB.Save(&round).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	message := "Motion published to users"
	if !input.IsMotionPublished {
		message = "Motion hidden from users"
	}
	c.JSON(http.StatusOK, gin.H{"message": message, "data": round})
}

// 5. LIHAT DAFTAR TURNAMEN (Baru)
func GetTournaments(c *gin.Context) {
	var tournaments []models.Tournament
	models.DB.Order("created_at desc").Find(&tournaments)
	c.JSON(http.StatusOK, gin.H{"data": tournaments})
}

// 5b. GET SINGLE TOURNAMENT BY ID
func GetTournament(c *gin.Context) {
	id := c.Param("id")
	var tournament models.Tournament
	if err := models.DB.First(&tournament, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tournament})
}

// 6. LIHAT DAFTAR TIM (Baru)
func GetTeams(c *gin.Context) {
	tournamentID := c.Query("tournament_id")
	var teams []models.Team
	query := models.DB.Preload("Speakers").Order("created_at desc")
	if tournamentID != "" {
		if _, err := strconv.Atoi(tournamentID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament_id"})
			return
		}
		query = query.Where("tournament_id = ?", tournamentID)
	}

	id := c.Query("id")
	if id != "" {
		if _, err := strconv.Atoi(id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}
		query = query.Where("id = ?", id)
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

// 6b. LIHAT DAFTAR SPEAKER (Berdasarkan Team ID)
func GetSpeakers(c *gin.Context) {
	teamID := c.Query("team_id")
	var speakers []models.Speaker
	query := models.DB.Preload("Team").Order("created_at asc")

	if teamID != "" {
		if _, err := strconv.Atoi(teamID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team_id"})
			return
		}
		query = query.Where("team_id = ?", teamID)
	}

	if err := query.Find(&speakers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if speakers == nil {
		speakers = []models.Speaker{}
	}
	c.JSON(http.StatusOK, gin.H{"data": speakers})
}

// 7. LIHAT DAFTAR RONDE (Berdasarkan Turnamen ID)
func GetRounds(c *gin.Context) {
	tournamentID := c.Query("tournament_id")
	var rounds []models.Round
	query := models.DB.Order("created_at asc")
	if tournamentID != "" {
		if _, err := strconv.Atoi(tournamentID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament_id"})
			return
		}
		query = query.Where("tournament_id = ?", tournamentID)
	}

	id := c.Query("id")
	if id != "" {
		if _, err := strconv.Atoi(id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}
		query = query.Where("id = ?", id)
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

// 7b. UPDATE ROUND STATUS (completed/in_progress)
func UpdateRoundStatus(c *gin.Context) {
	roundID := c.Param("id")
	var input struct {
		Status string `json:"status"` // "completed" or "in_progress"
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	if input.Status != "completed" && input.Status != "in_progress" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be 'completed' or 'in_progress'"})
		return
	}

	var round models.Round
	if err := models.DB.First(&round, roundID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Round not found"})
		return
	}

	round.Status = input.Status
	if err := models.DB.Save(&round).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": round, "message": "Round status updated"})
}

// 9. Buat Match (Pairing: Tim A vs Tim B)
func CreateMatch(c *gin.Context) {
	var input struct {
		RoundID       uint `json:"round_id"`
		GovTeamID     uint `json:"gov_team_id"`
		OppTeamID     uint `json:"opp_team_id"`
		RoomID        uint `json:"room_id"`
		AdjudicatorID uint `json:"adjudicator_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Debug log
	fmt.Printf("CreateMatch received: %+v\n", input)
	match := models.Match{
		RoundID:       input.RoundID,
		GovTeamID:     &input.GovTeamID,
		OppTeamID:     &input.OppTeamID,
		RoomID:        &input.RoomID,
		AdjudicatorID: &input.AdjudicatorID,
		IsCompleted:   false,
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
	roundID := c.Query("round_id")           // Filter per ronde
	tournamentID := c.Query("tournament_id") // Filter per tournament
	var matches []models.Match
	query := models.DB.Preload("GovTeam").Preload("OppTeam").Preload("Round").Preload("Room").Preload("Adjudicator").Order("id asc")

	if roundID != "" {
		if _, err := strconv.Atoi(roundID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid round_id"})
			return
		}
		query = query.Where("round_id = ?", roundID)
	}

	// Filter by tournament_id (join through rounds)
	if tournamentID != "" {
		if _, err := strconv.Atoi(tournamentID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament_id"})
			return
		}
		// Get all round IDs for this tournament first
		var roundIDs []uint
		models.DB.Model(&models.Round{}).Where("tournament_id = ?", tournamentID).Pluck("id", &roundIDs)
		if len(roundIDs) > 0 {
			query = query.Where("round_id IN ?", roundIDs)
		} else {
			// No rounds for this tournament, return empty
			c.JSON(http.StatusOK, gin.H{"data": []models.Match{}})
			return
		}
	}

	teamID := c.Query("team_id")
	if teamID != "" {
		if _, err := strconv.Atoi(teamID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team_id"})
			return
		}
		query = query.Where("gov_team_id = ? OR opp_team_id = ?", teamID, teamID)
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

// 11b. ASSIGN ADJUDICATOR PANEL TO MATCH
func AssignAdjudicatorPanel(c *gin.Context) {
	matchID := c.Param("id")
	var input struct {
		ChiefAdjID uint   `json:"chief_adj_id"`
		WingAdjIDs []uint `json:"wing_adj_ids"`
		PanelSize  int    `json:"panel_size"`
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

	// Set the chief adjudicator
	match.AdjudicatorID = &input.ChiefAdjID

	// Store wing adjudicators as comma-separated IDs in PanelJudges
	if len(input.WingAdjIDs) > 0 {
		wingIDsStr := ""
		for i, id := range input.WingAdjIDs {
			if i > 0 {
				wingIDsStr += ","
			}
			wingIDsStr += fmt.Sprintf("%d", id)
		}
		match.PanelJudges = wingIDsStr
	} else {
		match.PanelJudges = ""
	}

	if err := models.DB.Save(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload match with adjudicator data
	models.DB.Preload("Adjudicator").First(&match, matchID)

	c.JSON(http.StatusOK, gin.H{"data": match, "message": "Adjudicator panel assigned successfully"})
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
		if _, err := strconv.Atoi(tournamentID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament_id"})
			return
		}
		query = query.Where("tournament_id = ?", tournamentID)
	}
	if err := query.Find(&adjudicators).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if adjudicators == nil {
		adjudicators = []models.Adjudicator{}
	}

	// Create a response struct
	type AdjudicatorWithStats struct {
		models.Adjudicator
		AvgRating     float64 `json:"avg_rating"`
		MatchesJudged int64   `json:"matches_judged"`
	}

	var response []AdjudicatorWithStats

	for _, adj := range adjudicators {
		var avgRating float64
		// Calculate average rating
		models.DB.Model(&models.AdjudicatorFeedback{}).
			Where("adjudicator_id = ?", adj.ID).
			Select("COALESCE(AVG(rating), 0)").
			Scan(&avgRating)

		var matchesJudged int64
		// Count completed matches judged
		models.DB.Model(&models.Match{}).
			Where("adjudicator_id = ? AND is_completed = ?", adj.ID, true).
			Count(&matchesJudged)

		response = append(response, AdjudicatorWithStats{
			Adjudicator:   adj,
			AvgRating:     avgRating,
			MatchesJudged: matchesJudged,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
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
		if _, err := strconv.Atoi(tournamentID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament_id"})
			return
		}
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

// =====================================================
// CSV IMPORT FUNCTIONS
// =====================================================

// CSV Import Teams
// Format CSV: name,institution,speaker1,speaker2,speaker3
func ImportTeamsCSV(c *gin.Context) {
	tournamentID := c.Query("tournament_id")
	if tournamentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tournament_id is required"})
		return
	}

	tid, err := strconv.Atoi(tournamentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament_id"})
		return
	}

	var input struct {
		Data [][]string `json:"data"` // Array of rows, each row is array of columns
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := models.DB.Begin()
	teamsCreated := 0
	speakersCreated := 0

	for rowIdx, row := range input.Data {
		// Skip header row if detected
		if rowIdx == 0 && (row[0] == "name" || row[0] == "Name" || row[0] == "team" || row[0] == "Team") {
			continue
		}

		if len(row) < 2 {
			continue // Skip invalid rows
		}

		teamName := row[0]
		institution := row[1]

		if teamName == "" {
			continue
		}

		// Create team
		team := models.Team{
			TournamentID: uint(tid),
			Name:         teamName,
			Institution:  institution,
		}

		if err := tx.Create(&team).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create team '%s': %s", teamName, err.Error())})
			return
		}
		teamsCreated++

		// Create speakers (columns 3+)
		for i := 2; i < len(row); i++ {
			speakerName := row[i]
			if speakerName == "" {
				continue
			}

			speaker := models.Speaker{
				TeamID: team.ID,
				Name:   speakerName,
			}

			if err := tx.Create(&speaker).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create speaker '%s': %s", speakerName, err.Error())})
				return
			}
			speakersCreated++
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message":          "CSV imported successfully",
		"teams_created":    teamsCreated,
		"speakers_created": speakersCreated,
	})
}

// CSV Import Adjudicators
// Format CSV: name,institution
func ImportAdjudicatorsCSV(c *gin.Context) {
	tournamentID := c.Query("tournament_id")
	if tournamentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tournament_id is required"})
		return
	}

	tid, err := strconv.Atoi(tournamentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament_id"})
		return
	}

	var input struct {
		Data [][]string `json:"data"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := models.DB.Begin()
	created := 0

	for rowIdx, row := range input.Data {
		// Skip header row if detected
		if rowIdx == 0 && (row[0] == "name" || row[0] == "Name" || row[0] == "adjudicator" || row[0] == "Adjudicator") {
			continue
		}

		if len(row) < 1 || row[0] == "" {
			continue
		}

		adjName := row[0]
		institution := ""
		if len(row) > 1 {
			institution = row[1]
		}

		adj := models.Adjudicator{
			TournamentID: uint(tid),
			Name:         adjName,
			Institution:  institution,
		}

		if err := tx.Create(&adj).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create adjudicator '%s': %s", adjName, err.Error())})
			return
		}
		created++
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message":              "CSV imported successfully",
		"adjudicators_created": created,
	})
}

// CSV Import Rooms
// Format CSV: name,capacity (capacity is optional)
func ImportRoomsCSV(c *gin.Context) {
	tournamentID := c.Query("tournament_id")
	if tournamentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tournament_id is required"})
		return
	}

	tid, err := strconv.Atoi(tournamentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament_id"})
		return
	}

	var input struct {
		Data [][]string `json:"data"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := models.DB.Begin()
	created := 0

	for rowIdx, row := range input.Data {
		// Skip header row if detected
		if rowIdx == 0 && (row[0] == "name" || row[0] == "Name" || row[0] == "room" || row[0] == "Room") {
			continue
		}

		if len(row) < 1 || row[0] == "" {
			continue
		}

		roomName := row[0]
		capacity := 0
		if len(row) > 1 {
			cap, _ := strconv.Atoi(row[1])
			capacity = cap
		}

		room := models.Room{
			TournamentID: uint(tid),
			Name:         roomName,
			Capacity:     capacity,
		}

		if err := tx.Create(&room).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create room '%s': %s", roomName, err.Error())})
			return
		}
		created++
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message":       "CSV imported successfully",
		"rooms_created": created,
	})
}
