package models

import (
	"fmt"
	"log"
	"os" // <--- Jangan lupa ini

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Ambil URL dari Environment Variable (Render)
	dsn := os.Getenv("DATABASE_URL")

	// Kalau kosong, berarti kita lagi di Laptop (Localhost)
	if dsn == "" {
		dsn = "host=localhost user=admin password=rahasia123 dbname=eds_upi port=5433 sslmode=disable TimeZone=Asia/Jakarta"
		fmt.Println("💻 Mode: LOCALHOST (Docker)")
	} else {
		fmt.Println("☁️ Mode: CLOUD (Render)")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("❌ Gagal konek ke database!", err)
	}

	// Auto Migrate (Biar tabel otomatis dibuat di Supabase)
	err = database.AutoMigrate(
		&User{}, &Member{}, &Article{}, &CompetitionHistory{}, &Achievement{},
		&Tournament{}, &Team{}, &Speaker{}, &Round{}, &Match{}, &Ballot{},
		&Adjudicator{}, &Room{}, &AdjudicatorFeedback{},
	)

	DB = database
	fmt.Println("✅ SUKSES: Database Terhubung!")
}
