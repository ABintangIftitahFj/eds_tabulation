package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/star_fj/eds-backend/models"
	"golang.org/x/crypto/bcrypt"
)

var secretKey []byte

func init() {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	secretKey = []byte(key)
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 1. REGISTER (Buat Admin Baru)
func Register(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek Manual: Apakah username admin sudah ada?
	var existingUser models.User
	if err := models.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		// Kalau err == nil, artinya user KETEMU (Sudah ada)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username '" + input.Username + "' sudah dipakai! Coba login saja."})
		return
	}

	// Enkripsi Password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     "admin",
	}

	// Simpan ke Database
	if err := models.DB.Create(&user).Error; err != nil {
		// TAMPILKAN ERROR ASLI DI TERMINAL (Biar ketahuan kenapa gagal)
		fmt.Println("❌ ERROR SAAT CREATE USER:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Simpan ke DB: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin berhasil dibuat! Silakan Login."})
}

// 2. LOGIN (Cek Password & Dapat Token)
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Cari user berdasarkan username
	if err := models.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username tidak ditemukan!"})
		return
	}

	// Cek Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password Salah!"})
		return
	}

	// Generate Token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Expired 24 jam
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("failed to sign token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"role":  user.Role,
	})
}
