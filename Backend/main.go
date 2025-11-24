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
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 3. Setup Routing API
	api := r.Group("/api")
	{
		// ==============================
		// 🔐 AUTHENTICATION
		// ==============================
		api.POST("/register", controllers.Register) // Dev Only: Buat Admin
		api.POST("/login", controllers.Login)       // Login Admin

		// ==============================
		// 📰 FITUR COMPANY PROFILE
		// ==============================
		// Berita (Articles)
		api.GET("/articles", controllers.GetArticles)            // Public: Lihat semua berita
		api.GET("/articles/:slug", controllers.GetArticleBySlug) // Public: Baca 1 berita detail
		api.POST("/articles", controllers.CreateArticle)         // Admin: Tambah berita

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

		// Tim
		api.GET("/teams", controllers.GetTeams)    // <--- API untuk melihat daftar tim
		api.POST("/teams", controllers.CreateTeam) // <--- API untuk mendaftarkan tim baru

		// Ronde & Match

		// --- INPUT SKOR (TABULATOR) ---
		api.POST("/ballots", controllers.SubmitBallot)

		// RONDE
		api.GET("/rounds", controllers.GetRounds)
		api.POST("/rounds", controllers.CreateRound)

		// MATCHES
		api.GET("/matches", controllers.GetMatches)
		api.POST("/matches", controllers.CreateMatch)

		// STANDINGS (KLASEMEN)
		api.GET("/standings", controllers.GetStandings)
	}

	// 4. Jalankan Server
	r.Run(":8080")
}
