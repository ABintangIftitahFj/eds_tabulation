package models

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=admin password=password123 dbname=eds_upi port=5433 sslmode=disable TimeZone=Asia/Jakarta"

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

	DB = database
	fmt.Println("✅ SUKSES: Database Updated dengan Schema Baru!")
}
