package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/star_fj/eds-backend/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupControllerTestDB() {
	// Set required environment variables
	os.Setenv("JWT_SECRET", "test-secret-key-for-controllers")

	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// Run migrations
	models.DB = db
	models.DB.AutoMigrate(
		&models.Tournament{},
		&models.Team{},
		&models.Speaker{},
		&models.Round{},
		&models.Match{},
		&models.Ballot{},
		&models.Adjudicator{},
	)
}

func setupControllerTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := router.Group("/api")
	{
		// Tournament routes
		api.GET("/tournaments", GetTournaments)
		api.POST("/tournaments", CreateTournament)
		api.PUT("/tournaments/:id", UpdateTournament)
		api.DELETE("/tournaments/:id", DeleteTournament)

		// Team routes
		api.GET("/teams", GetTeams)
		api.POST("/teams", CreateTeam)
		api.DELETE("/teams/:id", DeleteTeam)

		// Round routes
		api.GET("/rounds", GetRounds)
		api.POST("/rounds", CreateRound)
		api.DELETE("/rounds/:id", DeleteRound)

		// Match routes
		api.GET("/matches", GetMatches)
		api.POST("/matches", CreateMatch)

		// Ballot routes
		api.POST("/submit-ballot", SubmitBallot)
		api.GET("/ballots", GetBallots)

		// Standings routes
		api.GET("/standings/teams", GetStandings)
		api.GET("/standings/speakers", GetSpeakerStandings)
		api.POST("/standings/recalculate", RecalculateStandings)
	}

	return router
}

func TestTournamentController(t *testing.T) {
	setupControllerTestDB()
	router := setupControllerTestRouter()

	t.Run("Create Tournament", func(t *testing.T) {
		tournament := models.Tournament{
			Name:        "Test Tournament",
			Slug:        "test-tournament",
			Format:      "asian",
			Location:    "Test Location",
			Description: "Test Description",
			Status:      "upcoming",
			IsPublic:    true,
		}

		jsonData, _ := json.Marshal(tournament)
		req, _ := http.NewRequest("POST", "/api/tournaments", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})

		assert.Equal(t, "Test Tournament", data["name"])
		assert.Equal(t, "test-tournament", data["slug"])
		assert.Equal(t, "asian", data["format"])
	})

	t.Run("Get Tournaments", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/tournaments", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].([]interface{})

		assert.GreaterOrEqual(t, len(data), 0)
	})
}

func TestTeamController(t *testing.T) {
	setupControllerTestDB()
	router := setupControllerTestRouter()

	// First create a tournament for testing
	tournament := models.Tournament{
		Name:   "Test Tournament",
		Status: "active",
	}
	models.DB.Create(&tournament)

	t.Run("Create Team", func(t *testing.T) {
		team := models.Team{
			Name:         "Test Team",
			Institution:  "Test University",
			TournamentID: tournament.ID,
		}

		jsonData, _ := json.Marshal(team)
		req, _ := http.NewRequest("POST", "/api/teams", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})

		assert.Equal(t, "Test Team", data["name"])
		assert.Equal(t, "Test University", data["institution"])
	})

	t.Run("Get Teams by Tournament", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/teams?tournament_id=1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].([]interface{})

		assert.GreaterOrEqual(t, len(data), 0)
	})
}

func TestBallotController(t *testing.T) {
	setupControllerTestDB()
	router := setupControllerTestRouter()

	// Setup test data
	tournament := models.Tournament{Name: "Test Tournament", Status: "active"}
	models.DB.Create(&tournament)

	govTeam := models.Team{Name: "Gov Team", TournamentID: tournament.ID}
	oppTeam := models.Team{Name: "Opp Team", TournamentID: tournament.ID}
	models.DB.Create(&govTeam)
	models.DB.Create(&oppTeam)

	round := models.Round{Name: "Round 1", TournamentID: tournament.ID}
	models.DB.Create(&round)

	match := models.Match{
		RoundID:   round.ID,
		GovTeamID: &govTeam.ID,
		OppTeamID: &oppTeam.ID,
	}
	models.DB.Create(&match)

	adjudicator := models.Adjudicator{Name: "Test Judge", TournamentID: tournament.ID}
	models.DB.Create(&adjudicator)

	t.Run("Submit Ballot", func(t *testing.T) {
		ballotData := map[string]interface{}{
			"match_id":       match.ID,
			"adjudicator_id": adjudicator.ID,
			"winner":         "gov",
			"scores": []map[string]interface{}{
				{
					"speaker":   map[string]string{"name": "Gov PM"},
					"score":     85,
					"position":  "PM",
					"team_role": "gov",
				},
				{
					"speaker":   map[string]string{"name": "Gov DPM"},
					"score":     80,
					"position":  "DPM",
					"team_role": "gov",
				},
				{
					"speaker":   map[string]string{"name": "Opp LO"},
					"score":     78,
					"position":  "LO",
					"team_role": "opp",
				},
				{
					"speaker":   map[string]string{"name": "Opp DLO"},
					"score":     82,
					"position":  "DLO",
					"team_role": "opp",
				},
			},
		}

		jsonData, _ := json.Marshal(ballotData)
		req, _ := http.NewRequest("POST", "/api/submit-ballot", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response["message"], "Skor disimpan")
		assert.Equal(t, float64(govTeam.ID), response["winner_id"])
		assert.Equal(t, float64(165), response["total_gov"]) // 85 + 80
		assert.Equal(t, float64(160), response["total_opp"]) // 78 + 82
	})

	t.Run("Get Ballots by Match", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/ballots?match_id=1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].([]interface{})

		assert.Equal(t, 4, len(data)) // 4 speakers
	})
}

func TestErrorHandling(t *testing.T) {
	setupControllerTestDB()
	router := setupControllerTestRouter()

	t.Run("Tournament Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/tournaments/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "not found")
	})

	t.Run("Invalid JSON Input", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/tournaments", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
