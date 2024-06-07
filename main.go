package main

import (
	"gin/middleware"
	"gin/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	routes.SetupRouter()

	router.Use(middleware.RateLimitMiddleware(50, 1*60*60*1000000000)) // 50 requests per hour

	// Player stats routes
	router.POST("/player-stats", middleware.AuthMiddleware(), routes.Add)
	router.POST("/batch-insert", routes.BatchInsertCsvData)
	router.GET("/player-stats/:id", middleware.AuthMiddleware(), routes.GetById)
	router.GET("/player-stats/player/:name/:surname", middleware.AuthMiddleware(), routes.GetByNameAndSurname)
	router.GET("/player-stats", middleware.AuthMiddleware(), routes.GetAll)
	router.PUT("/player-stats/:id", middleware.AuthMiddleware(), routes.Update)
	router.DELETE("/player-stats/:id", middleware.AuthMiddleware(), routes.Delete)

	// Register routes
	router.POST("/register", middleware.RegisterHandler)
	router.POST("/login", middleware.LoginHandler)
	router.POST("/change-password", middleware.AuthMiddleware(), middleware.ChangePasswordHandler)

	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
