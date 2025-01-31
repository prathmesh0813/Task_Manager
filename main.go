package main

import (
	"task_manager/dao"
	"task_manager/logger"
	"task_manager/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("", "failed to load .env file", err.Error())
	}

	logger.InitLogger()
	defer logger.InitLogger()

	logger.Info("", "Starting the application")

	dao.InitDB()
	logger.Info("", "Database connection initialized")

	server := gin.Default()
	logger.Info("", "Server initialized successfully")

	server.Use(cors.Default())

	routes.RegisterRoutes(server)
	logger.Info("", "Routes registered successfully")

	if err := server.Run(":8080"); err != nil {
		logger.Error("", "failed to start the server", err.Error())
	}
}
