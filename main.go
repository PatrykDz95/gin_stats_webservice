package main

import (
	"gin/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("/player_stats", func(c *gin.Context) {
		stats, err := service.GetPlayerStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stats)
	})

	router.Run("localhost:8080")
}
