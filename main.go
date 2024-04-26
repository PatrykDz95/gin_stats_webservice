package main

import (
	routes "gin/routers"
	"gin/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("/player_stats", func(c *gin.Context) {
		stats, err := service.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stats)
	})
	routes.SetupRouter()
	router.POST("/player-stats", routes.InsertPlayerStats)

	router.Run("localhost:8080")
}
