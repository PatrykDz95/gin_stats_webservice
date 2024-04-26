package routes

import (
	"gin/database"
	_ "gin/database"
	"gin/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Initialize database connection
	database.InitDB()

	// Route to handle inserting data

	return router
}

func InsertPlayerStats(c *gin.Context) {
	playerStats := &models.PlayerStats{}
	if err := c.ShouldBindJSON(playerStats); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created := database.DB.Create(playerStats)

	c.JSON(http.StatusCreated, gin.H{"player": created.Statement.Model})
}
