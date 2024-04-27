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
	router.POST("/player-stats", routes.Add)
	router.GET("/player-stats/:id", routes.Get)
	router.GET("/player-stats", routes.GetAll)
	router.PUT("/player-stats/:id", routes.Update)

	router.Run("localhost:8080")
}
