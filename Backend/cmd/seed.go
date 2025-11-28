package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/star_fj/eds-backend/models"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using existing environment variables")
	}

	// 1. Konek Database
	models.ConnectDatabase()

	// 2. Seed Admin User
	seedAdmin()

	// 3. Seed Tournament Data
	seedTournamentData()
}

func seedAdmin() {
	username := os.Getenv("SEED_ADMIN_USERNAME")
	if username == "" {
		username = "admin"
	}

	password := os.Getenv("SEED_ADMIN_PASSWORD")
	if password == "" {
		password = "admin123"
	}

	var existingUser models.User
	if err := models.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		fmt.Println("⚠️  User 'admin' sudah ada!")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("❌ Gagal mengenkripsi password:", err)
	}
	user := models.User{
		Username: username,
		Password: string(hashedPassword),
		Role:     "admin",
	}

	if err := models.DB.Create(&user).Error; err != nil {
		log.Fatal("❌ Gagal membuat user:", err)
	}

	fmt.Println("✅ SUKSES! User admin berhasil dibuat.")
	fmt.Println("👉 Username: " + username)
	fmt.Println("👉 Password: " + password)
}

func seedTournamentData() {
	fmt.Println("🏆 Mulai seeding tournament data...")

	// 1. Buat Tournament
	tournament := models.Tournament{
		Name:        "EDS Championship 2025",
		Slug:        "eds-championship-2025",
		Format:      "asian",
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 0, 3),
		Location:    "Universitas Pendidikan Indonesia",
		Description: "Tournament debat terbesar tahun ini dengan format Asian Parliamentary",
		Status:      "upcoming",
		IsPublic:    true,
	}

	if err := models.DB.Create(&tournament).Error; err != nil {
		log.Fatal("❌ Gagal membuat tournament:", err)
	}
	fmt.Printf("✅ Tournament '%s' berhasil dibuat (ID: %d)\n", tournament.Name, tournament.ID)

	// 2. Buat 16 Tim
	teams := []models.Team{
		{TournamentID: tournament.ID, Name: "UGM A", Institution: "Universitas Gadjah Mada"},
		{TournamentID: tournament.ID, Name: "UI A", Institution: "Universitas Indonesia"},
		{TournamentID: tournament.ID, Name: "ITB A", Institution: "Institut Teknologi Bandung"},
		{TournamentID: tournament.ID, Name: "UNPAD A", Institution: "Universitas Padjadjaran"},
		{TournamentID: tournament.ID, Name: "UPI A", Institution: "Universitas Pendidikan Indonesia"},
		{TournamentID: tournament.ID, Name: "UNDIP A", Institution: "Universitas Diponegoro"},
		{TournamentID: tournament.ID, Name: "UNAIR A", Institution: "Universitas Airlangga"},
		{TournamentID: tournament.ID, Name: "UB A", Institution: "Universitas Brawijaya"},
		{TournamentID: tournament.ID, Name: "USU A", Institution: "Universitas Sumatera Utara"},
		{TournamentID: tournament.ID, Name: "UNHAS A", Institution: "Universitas Hasanuddin"},
		{TournamentID: tournament.ID, Name: "UGM B", Institution: "Universitas Gadjah Mada"},
		{TournamentID: tournament.ID, Name: "UI B", Institution: "Universitas Indonesia"},
		{TournamentID: tournament.ID, Name: "ITB B", Institution: "Institut Teknologi Bandung"},
		{TournamentID: tournament.ID, Name: "UNPAD B", Institution: "Universitas Padjadjaran"},
		{TournamentID: tournament.ID, Name: "UPI B", Institution: "Universitas Pendidikan Indonesia"},
		{TournamentID: tournament.ID, Name: "UNDIP B", Institution: "Universitas Diponegoro"},
	}

	for _, team := range teams {
		if err := models.DB.Create(&team).Error; err != nil {
			log.Printf("❌ Gagal membuat tim %s: %v", team.Name, err)
			continue
		}

		// Buat speakers untuk setiap tim (3 speakers per tim untuk Asian Parliamentary)
		speakers := []models.Speaker{
			{TeamID: team.ID, Name: fmt.Sprintf("%s Speaker 1", team.Name)},
			{TeamID: team.ID, Name: fmt.Sprintf("%s Speaker 2", team.Name)},
			{TeamID: team.ID, Name: fmt.Sprintf("%s Speaker 3", team.Name)},
		}

		for _, speaker := range speakers {
			if err := models.DB.Create(&speaker).Error; err != nil {
				log.Printf("❌ Gagal membuat speaker %s: %v", speaker.Name, err)
			}
		}
	}
	fmt.Printf("✅ 16 tim berhasil dibuat dengan masing-masing 3 speakers\n")

	// 3. Buat Adjudicators (Juri)
	adjudicators := []models.Adjudicator{
		{TournamentID: tournament.ID, Name: "Prof. Dr. Ahmad Santoso", Institution: "UI", Level: "Chief"},
		{TournamentID: tournament.ID, Name: "Dr. Siti Nurhaliza", Institution: "UGM", Level: "Chief"},
		{TournamentID: tournament.ID, Name: "Dr. Budi Pratama", Institution: "ITB", Level: "Wing"},
		{TournamentID: tournament.ID, Name: "Dr. Lisa Wahyuni", Institution: "UNPAD", Level: "Wing"},
		{TournamentID: tournament.ID, Name: "M.A. Reza Fahlevi", Institution: "UPI", Level: "Wing"},
		{TournamentID: tournament.ID, Name: "M.A. Diana Sari", Institution: "UNDIP", Level: "Wing"},
		{TournamentID: tournament.ID, Name: "Andi Surya Permana", Institution: "UNAIR", Level: "Panelist"},
		{TournamentID: tournament.ID, Name: "Maya Indrawati", Institution: "UB", Level: "Panelist"},
		{TournamentID: tournament.ID, Name: "Hendro Wijaya", Institution: "USU", Level: "Panelist"},
		{TournamentID: tournament.ID, Name: "Dewi Kartika", Institution: "UNHAS", Level: "Panelist"},
		{TournamentID: tournament.ID, Name: "Fajar Ramadhan", Institution: "UGM", Level: "Panelist"},
		{TournamentID: tournament.ID, Name: "Nova Safitri", Institution: "UI", Level: "Panelist"},
	}

	for _, adj := range adjudicators {
		if err := models.DB.Create(&adj).Error; err != nil {
			log.Printf("❌ Gagal membuat adjudicator %s: %v", adj.Name, err)
		}
	}
	fmt.Printf("✅ %d adjudicators berhasil dibuat\n", len(adjudicators))

	// 4. Buat Rooms
	rooms := []models.Room{
		{TournamentID: tournament.ID, Name: "A1", Location: "Gedung A Lantai 1", Capacity: 30},
		{TournamentID: tournament.ID, Name: "A2", Location: "Gedung A Lantai 1", Capacity: 30},
		{TournamentID: tournament.ID, Name: "A3", Location: "Gedung A Lantai 1", Capacity: 30},
		{TournamentID: tournament.ID, Name: "B1", Location: "Gedung B Lantai 1", Capacity: 35},
		{TournamentID: tournament.ID, Name: "B2", Location: "Gedung B Lantai 1", Capacity: 35},
		{TournamentID: tournament.ID, Name: "B3", Location: "Gedung B Lantai 1", Capacity: 35},
		{TournamentID: tournament.ID, Name: "C1", Location: "Gedung C Lantai 1", Capacity: 40},
		{TournamentID: tournament.ID, Name: "C2", Location: "Gedung C Lantai 1", Capacity: 40},
		{TournamentID: tournament.ID, Name: "D1", Location: "Gedung D Lantai 1", Capacity: 25},
		{TournamentID: tournament.ID, Name: "D2", Location: "Gedung D Lantai 1", Capacity: 25},
		{TournamentID: tournament.ID, Name: "E1", Location: "Gedung E Lantai 2", Capacity: 30},
		{TournamentID: tournament.ID, Name: "E2", Location: "Gedung E Lantai 2", Capacity: 30},
	}

	for _, room := range rooms {
		if err := models.DB.Create(&room).Error; err != nil {
			log.Printf("❌ Gagal membuat room %s: %v", room.Name, err)
		}
	}
	fmt.Printf("✅ %d rooms berhasil dibuat\n", len(rooms))

	// 5. Buat Rounds
	rounds := []models.Round{
		{TournamentID: tournament.ID, Name: "Penyisihan 1", Motion: "THW ban social media for individuals under 18", IsPublished: true},
		{TournamentID: tournament.ID, Name: "Penyisihan 2", Motion: "THW implement universal basic income", IsPublished: true},
		{TournamentID: tournament.ID, Name: "Penyisihan 3", Motion: "THW prioritize environmental protection over economic growth", IsPublished: true},
		{TournamentID: tournament.ID, Name: "Quarter Final", Motion: "THW abolish the death penalty worldwide", IsPublished: false},
		{TournamentID: tournament.ID, Name: "Semi Final", Motion: "THW make voting mandatory in democratic elections", IsPublished: false},
		{TournamentID: tournament.ID, Name: "Final", Motion: "THW establish a world government", IsPublished: false},
	}

	for _, round := range rounds {
		if err := models.DB.Create(&round).Error; err != nil {
			log.Printf("❌ Gagal membuat round %s: %v", round.Name, err)
		}
	}
	fmt.Printf("✅ %d rounds berhasil dibuat\n", len(rounds))

	fmt.Println("🎉 SELESAI! Semua seed data berhasil dibuat:")
	fmt.Printf("   - 1 Tournament: %s\n", tournament.Name)
	fmt.Printf("   - 16 Teams dengan masing-masing 3 speakers\n")
	fmt.Printf("   - %d Adjudicators\n", len(adjudicators))
	fmt.Printf("   - %d Rooms\n", len(rooms))
	fmt.Printf("   - %d Rounds\n", len(rounds))
}
