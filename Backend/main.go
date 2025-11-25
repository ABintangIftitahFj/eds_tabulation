package main

import (
	"github.com/gin-gonic/gin"
	"github.com/star_fj/eds-backend/controllers"
	"github.com/star_fj/eds-backend/models"
)

func main() {
	// 1. Nyalakan Koneksi Database (Auto Migrate tabel baru)
	models.ConnectDatabase()

	// 2. Siapkan Server Gin
	r := gin.Default()

	// Tambahkan CORS Middleware agar Frontend bisa akses Backend
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 3. Setup Routing API
	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		api.POST("/upload", controllers.UploadFile)

		// ==============================
		// 🔐 AUTHENTICATION
		// ==============================
		api.POST("/register", controllers.Register) // Dev Only: Buat Admin
		api.POST("/login", controllers.Login)       // Login Admin

		// ==============================
		// 📰 FITUR COMPANY PROFILE
		// ==============================
		// Berita (Articles)
		api.GET("/articles", controllers.GetArticles)                   // Public: Lihat semua berita
		api.POST("/articles", controllers.CreateArticle)                // Admin: Tambah berita
		api.GET("/articles/detail/:slug", controllers.GetArticleBySlug) // Public: Baca 1 berita detail
		api.PUT("/articles/:id", controllers.UpdateArticle)             // Admin: Update berita
		api.DELETE("/articles/:id", controllers.DeleteArticle)          // Admin: Hapus berita

		// Nanti tambah route Member di sini:
		// api.GET("/members", controllers.GetMembers)

		// ==============================
		// 🏆 FITUR TABULASI (TABBYCAT)
		// ==============================
		// Nanti tambah route Submit Ballot di sini:
		// api.POST("/submit-ballot", controllers.SubmitBallot)

		// Turnamen
		api.GET("/tournaments", controllers.GetTournaments)
		api.POST("/tournaments", controllers.CreateTournament)
		api.PUT("/tournaments/:id", controllers.UpdateTournament)
		api.DELETE("/tournaments/:id", controllers.DeleteTournament)

		// Tim
		api.GET("/teams", controllers.GetTeams)    // <--- API untuk melihat daftar tim
		api.POST("/teams", controllers.CreateTeam) // <--- API untuk mendaftarkan tim baru
		api.DELETE("/teams/:id", controllers.DeleteTeam)

		// Ronde & Match

		// --- INPUT SKOR (TABULATOR) ---
		api.POST("/ballots", controllers.SubmitBallot)
		api.GET("/ballots", controllers.GetBallots)

		// RONDE
		api.GET("/rounds", controllers.GetRounds)
		api.POST("/rounds", controllers.CreateRound)
		api.DELETE("/rounds/:id", controllers.DeleteRound)
		api.PUT("/rounds/:id/publish-draw", controllers.PublishDraw)
		api.PUT("/rounds/:id/publish-motion", controllers.PublishMotion)

		// MATCHES
		api.GET("/matches", controllers.GetMatches)
		api.POST("/matches", controllers.CreateMatch)
		api.PUT("/matches/:id/result", controllers.UpdateMatchResult)
		api.DELETE("/matches/:id", controllers.DeleteMatch)

		// ADJUDICATORS
		api.GET("/adjudicators", controllers.GetAdjudicators)
		api.POST("/adjudicators", controllers.CreateAdjudicator)
		api.DELETE("/adjudicators/:id", controllers.DeleteAdjudicator)

		// ROOMS
		api.GET("/rooms", controllers.GetRooms)
		api.POST("/rooms", controllers.CreateRoom)
		api.DELETE("/rooms/:id", controllers.DeleteRoom)

		// STANDINGS (KLASEMEN)
		api.GET("/standings", controllers.GetStandings) // Legacy support if needed
		api.GET("/standings/teams", controllers.GetStandings)
		api.GET("/standings/speakers", controllers.GetSpeakerStandings)

		// INSTITUTIONS
		api.GET("/institutions", controllers.GetParticipatingInstitutions)

		// ADJUDICATOR FEEDBACK (USER RATING)
		api.GET("/adjudicator-feedback/check", func(c *gin.Context) {
			controllers.CheckFeedbackExists(c, models.DB)
		})
		api.GET("/adjudicator-feedback", func(c *gin.Context) {
			controllers.GetAdjudicatorFeedback(c, models.DB)
		})
		api.POST("/adjudicator-feedback", func(c *gin.Context) {
			controllers.CreateAdjudicatorFeedback(c, models.DB)
		})
		api.GET("/adjudicator-feedback/stats/:adjudicator_id", func(c *gin.Context) {
			controllers.GetFeedbackStats(c, models.DB)
		})
		api.DELETE("/adjudicator-feedback/:id", func(c *gin.Context) {
			controllers.DeleteAdjudicatorFeedback(c, models.DB)
		})
	}

	// 4. Jalankan Server
	r.Run(":8080")
}
