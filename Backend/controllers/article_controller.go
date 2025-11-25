package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

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

// generateSlug creates a URL-friendly slug from the title
func generateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special chars except hyphens and alphanumeric
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")
	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")
	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	// If slug is empty, generate a random one
	if slug == "" {
		slug = fmt.Sprintf("article-%d", time.Now().UnixNano())
	}

	return slug
}

// ensureUniqueSlug makes sure slug is unique in database
func ensureUniqueSlug(slug string) string {
	var count int64
	baseSlug := slug
	counter := 1

	for {
		models.DB.Model(&models.Article{}).Where("slug = ?", slug).Count(&count)
		if count == 0 {
			break
		}
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}

	return slug
}

// 2. POST Article (Untuk Admin nambah berita)
func CreateArticle(c *gin.Context) {
	var input models.Article

	// Validasi Input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate slug dari title dan pastikan unik
	input.Slug = generateSlug(input.Title)
	input.Slug = ensureUniqueSlug(input.Slug)

	// Set publish date jika belum ada
	if input.PublishDate.IsZero() {
		input.PublishDate = time.Now()
	}

	// Simpan ke Database
	if err := models.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article: " + err.Error()})
		return
	}

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

// 4. PUT Article (Update berita)
func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article

	// Cari artikel berdasarkan ID
	if err := models.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artikel tidak ditemukan"})
		return
	}

	var input models.Article
	// Validasi Input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate slug baru jika title berubah
	if input.Title != article.Title && input.Title != "" {
		input.Slug = generateSlug(input.Title)
	}

	// Update artikel
	if err := models.DB.Model(&article).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update artikel: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": article})
}

// 5. DELETE Article (Hapus berita)
func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article

	// Cari artikel berdasarkan ID
	if err := models.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artikel tidak ditemukan"})
		return
	}

	// Hapus artikel
	if err := models.DB.Delete(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus artikel: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artikel berhasil dihapus"})
}
