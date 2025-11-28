package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB_LOCALHOST *gorm.DB

// ConnectDatabaseLocalhost - Kode lama untuk localhost saja
func ConnectDatabaseLocalhost() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("❌ Gagal konek ke database!", err)
	}

	// AUTO MIGRATE: Daftarkan SEMUA Struct baru di sini
	err = database.AutoMigrate(
		&User{},
		// Company Profile
		&Member{},
		&Article{},
		&CompetitionHistory{}, // <-- Baru
		&Achievement{},
		// Tabulation System
		&Tournament{},  // <-- Baru
		&Adjudicator{}, // <-- Juri
		&Room{},        // <-- Ruangan
		&Team{},
		&Speaker{},
		&Round{},
		&Match{},
		&Ballot{},
		&AdjudicatorFeedback{}, // <-- Feedback Juri
		// Motion (opsional jika dipisah)
	)

	if err != nil {
		log.Fatal("❌ Gagal migrasi tabel:", err)
	}

	DB_LOCALHOST = database
	fmt.Println("✅ SUKSES: Database Localhost Terhubung!")
}
