package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/star_fj/eds-backend/controllers"
	"github.com/star_fj/eds-backend/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	// Set environment variables for testing
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_PORT", "3306")

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

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := router.Group("/api")
	{
		// Tournament routes
		api.GET("/tournaments", controllers.GetTournaments)
		api.POST("/tournaments", controllers.CreateTournament)
		api.PUT("/tournaments/:id", controllers.UpdateTournament)
		api.DELETE("/tournaments/:id", controllers.DeleteTournament)

		// Team routes
		api.GET("/teams", controllers.GetTeams)
		api.POST("/teams", controllers.CreateTeam)
		api.DELETE("/teams/:id", controllers.DeleteTeam)

		// Round routes
		api.GET("/rounds", controllers.GetRounds)
		api.POST("/rounds", controllers.CreateRound)
		api.DELETE("/rounds/:id", controllers.DeleteRound)

		// Match routes
		api.GET("/matches", controllers.GetMatches)
		api.POST("/matches", controllers.CreateMatch)

		// Ballot routes
		api.POST("/submit-ballot", controllers.SubmitBallot)
		api.GET("/ballots", controllers.GetBallots)

		// Standings routes
		api.GET("/standings/teams", controllers.GetStandings)
		api.GET("/standings/speakers", controllers.GetSpeakerStandings)
		api.POST("/standings/recalculate", controllers.RecalculateStandings)
	}

	return router
}

func TestTournamentAPI(t *testing.T) {
	setupTestDB()
	router := setupTestRouter()

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

		assert.Greater(t, len(data), 0)
	})

	t.Run("Update Tournament", func(t *testing.T) {
		// First create a tournament
		tournament := models.Tournament{
			Name:   "Tournament to Update",
			Status: "upcoming",
		}
		models.DB.Create(&tournament)

		updateData := models.Tournament{
			Status: "active",
		}

		jsonData, _ := json.Marshal(updateData)
		req, _ := http.NewRequest("PUT", "/api/tournaments/1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})

		assert.Equal(t, "active", data["status"])
	})
}

func TestTeamAPI(t *testing.T) {
	setupTestDB()
	router := setupTestRouter()

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

		assert.Greater(t, len(data), 0)
	})
}

func TestBallotAPI(t *testing.T) {
	setupTestDB()
	router := setupTestRouter()

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

func TestStandingsAPI(t *testing.T) {
	setupTestDB()
	router := setupTestRouter()

	// Setup test data with completed matches
	tournament := models.Tournament{Name: "Test Tournament", Status: "active"}
	models.DB.Create(&tournament)

	team1 := models.Team{Name: "Team 1", TournamentID: tournament.ID, TotalVP: 2, TotalSpeaker: 500, Wins: 2, Losses: 0}
	team2 := models.Team{Name: "Team 2", TournamentID: tournament.ID, TotalVP: 1, TotalSpeaker: 450, Wins: 1, Losses: 1}
	models.DB.Create(&team1)
	models.DB.Create(&team2)

	t.Run("Get Team Standings", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/standings/teams?tournament_id=1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].([]interface{})

		assert.Equal(t, 2, len(data))

		// Check if sorted by VP and speaker points
		firstTeam := data[0].(map[string]interface{})
		assert.Equal(t, "Team 1", firstTeam["name"])
		assert.Equal(t, float64(2), firstTeam["total_vp"])
	})

	t.Run("Get Speaker Standings", func(t *testing.T) {
		// Add some speakers with scores
		speaker1 := models.Speaker{Name: "Speaker 1", TeamID: team1.ID, TotalScore: 255}
		speaker2 := models.Speaker{Name: "Speaker 2", TeamID: team2.ID, TotalScore: 240}
		models.DB.Create(&speaker1)
		models.DB.Create(&speaker2)

		req, _ := http.NewRequest("GET", "/api/standings/speakers?tournament_id=1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].([]interface{})

		assert.Greater(t, len(data), 0)
	})
}

func TestErrorHandling(t *testing.T) {
	setupTestDB()
	router := setupTestRouter()

	t.Run("Tournament Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/tournaments/999", nil)
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

	t.Run("Missing Required Fields", func(t *testing.T) {
		incompleteTeam := map[string]interface{}{
			"name": "", // Empty required field
		}

		jsonData, _ := json.Marshal(incompleteTeam)
		req, _ := http.NewRequest("POST", "/api/teams", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should handle validation errors gracefully
		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

func BenchmarkTournamentCreation(b *testing.B) {
	setupTestDB()
	router := setupTestRouter()

	tournament := models.Tournament{
		Name:        "Benchmark Tournament",
		Slug:        "benchmark-tournament",
		Format:      "asian",
		Location:    "Benchmark Location",
		Description: "Benchmark Description",
		Status:      "upcoming",
		IsPublic:    true,
	}

	jsonData, _ := json.Marshal(tournament)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/api/tournaments", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("Expected status OK, got %d", w.Code)
		}
	}
}
