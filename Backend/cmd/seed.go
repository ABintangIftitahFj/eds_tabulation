package main

import (
	"fmt"
	"log"

	"github.com/star_fj/eds-backend/models"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 1. Konek Database
	models.ConnectDatabase()

	// 2. Data Admin yang mau dibuat
	username := "admin"
	password := "admin123"

	// 3. Cek apakah sudah ada?
	var existingUser models.User
	if err := models.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		fmt.Println("⚠️  User 'admin' sudah ada!")

		// Opsional: Reset password kalau mau
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		existingUser.Password = string(hashedPassword)
		models.DB.Save(&existingUser)
		fmt.Println("✅ Password untuk 'admin' berhasil di-RESET menjadi: " + password)
		return
	}

	// 4. Buat Baru
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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
