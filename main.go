package main

import (
	"task_manager/dao"
	"task_manager/routes"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"github.com/gin-contrib/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		utils.Logger.Fatal("Error loading .env file")
	}

	utils.InitLogger()
	defer utils.InitLogger()

	utils.Logger.Info("Starting the application...")

	dao.InitDB()
	utils.Logger.Info("Database connection initialized")

	server := gin.Default()
	utils.Logger.Info("Server initialized")

	server.Use(cors.Default())

	routes.RegisterRoutes(server)
	utils.Logger.Info("Routes registered")

	if err := server.Run(":8080"); err != nil {
		utils.Logger.Fatal("Failed to start the server", zap.Error(err))
	}
}
