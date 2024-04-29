package routes

import (
	"gin/database"
	_ "gin/database"
	"gin/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
	playerStats.CreatedOn = time.Now().Format("2006-01-02 15:04:05")
	if err := c.ShouldBindJSON(&playerStats); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	created := database.DB.Create(&playerStats)
	if created.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"errorMessage": created.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"player": created.Statement.Model})
}

func GetAll(c *gin.Context) {
	var playerStats []*models.PlayerStats
	found := database.DB.Find(&playerStats)
	if found.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"errorMessage": found.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"players": found.Statement.Model})
}

func Get(c *gin.Context) {
	playerStats := models.PlayerStats{}
	found := database.DB.First(&playerStats, c.Param("id"))
	if found.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"errorMessage": found.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"player": found.Statement.Model})

}

func Update(c *gin.Context) {
	playerStats := models.PlayerStats{}
	if err := c.ShouldBindJSON(&playerStats); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	updated := database.DB.Model(&models.PlayerStats{}).Where("id = ?", c.Param("id")).Updates(playerStats)
	if updated.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"errorMessage": updated.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"player": updated.Statement.Model})
}

func Delete(c *gin.Context) {
	playerStats := models.PlayerStats{}
	if err := c.ShouldBindJSON(&playerStats); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	deleted := database.DB.Delete(&playerStats, c.Param("id"))
	if deleted.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"errorMessage": deleted.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"player": deleted.Statement.Model})
}
