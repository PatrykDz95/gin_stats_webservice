package main

import (
	"gin/models"
	routes "gin/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	routes.SetupRouter()

	router.POST("/player-stats", models.AuthMiddleware(), routes.Add)
	router.GET("/player-stats/:id", models.AuthMiddleware(), routes.GetById)
	router.GET("/player-stats/player/:name/:surname", models.AuthMiddleware(), routes.GetByNameAndSurname)
	router.GET("/player-stats", models.AuthMiddleware(), routes.GetAll)
	router.PUT("/player-stats/:id", models.AuthMiddleware(), routes.Update)
	router.DELETE("/player-stats/:id", models.AuthMiddleware(), routes.Delete)

	// Register routes
	router.POST("/register", models.RegisterHandler)
	router.POST("/login", models.LoginHandler)
	router.POST("/change-password", models.AuthMiddleware(), models.ChangePasswordHandler)

	router.Run("localhost:8080")
}
