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

func Add(c *gin.Context) {
	playerStats := models.PlayerStats{}
	if err := c.ShouldBindJSON(playerStats); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created := database.DB.Create(playerStats)

	c.JSON(http.StatusCreated, gin.H{"player": created.Statement.Model})
}

func GetAll(c *gin.Context) {
	var playerStats []*models.PlayerStats
	found := database.DB.Find(&playerStats)

	c.JSON(http.StatusCreated, gin.H{"players": found.Statement.Model})
}

func Get(c *gin.Context) {
	playerStats := models.PlayerStats{}
	found := database.DB.First(&playerStats, c.Param("id"))

	c.JSON(http.StatusCreated, gin.H{"players": found.Statement.Model})
}
