package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star_fj/eds-backend/models"
)

// 1. GET All Articles (Untuk Halaman News)
func GetArticles(c *gin.Context) {
	var articles []models.Article
	// Ambil semua data dari database
	models.DB.Find(&articles)

	c.JSON(http.StatusOK, gin.H{"data": articles})
}

// 2. POST Article (Untuk Admin nambah berita)
func CreateArticle(c *gin.Context) {
	var input models.Article

	// Validasi Input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simpan ke Database
	models.DB.Create(&input)

	c.JSON(http.StatusOK, gin.H{"data": input})
}

// 3. GET Single Article (Detail Berita)
func GetArticleBySlug(c *gin.Context) {
	slug := c.Param("slug")
	var article models.Article

	// Cari berdasarkan slug (misal: /api/news/juara-1-nudc)
	if err := models.DB.Where("slug = ?", slug).First(&article).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": article})
}
