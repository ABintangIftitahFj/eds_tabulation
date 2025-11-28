package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set environment variables for testing
	os.Setenv("JWT_SECRET", "test-secret-key")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Simple test endpoint to check what tournament IDs exist
	r.GET("/test/tournaments", func(c *gin.Context) {
		// Return mock data for testing
		tournaments := []map[string]interface{}{
			{"id": 1, "name": "PIMNAS 37", "status": "completed"},
			{"id": 2, "name": "Test Tournament", "status": "active"},
		}
		c.JSON(http.StatusOK, gin.H{"data": tournaments})
	})

	r.GET("/test/tournaments/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Mock tournament data
		tournaments := map[string]interface{}{
			"1": map[string]interface{}{"id": 1, "name": "PIMNAS 37", "status": "completed"},
			"2": map[string]interface{}{"id": 2, "name": "Test Tournament", "status": "active"},
		}

		if tournament, exists := tournaments[id]; exists {
			c.JSON(http.StatusOK, gin.H{"data": tournament})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Tournament with ID %s not found", id)})
		}
	})

	log.Println("Test server running on http://localhost:8080")
	r.Run(":8080")
}
