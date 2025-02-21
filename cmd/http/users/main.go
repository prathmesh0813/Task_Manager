package main

import (
	"path/filepath"
	"task_manager/dao"
	"task_manager/logger"
	"task_manager/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	
	logger.InitLogger()
	defer logger.InitLogger()

	rootDir,err :=filepath.Abs("../../../")
	if err!=nil{
		logger.Error("","error getting root directory", err.Error())
	}

	envPath := filepath.Join(rootDir,".env")
	err = godotenv.Load(envPath)
	if err!=nil{
		logger.Error("","failed to load .env file", err.Error())
	}

	logger.Info("", "Starting the application")

	dao.InitDB()
	logger.Info("", "Database connection initialized")

	server := gin.Default()
	logger.Info("", "Server initialized successfully")

	server.Use(cors.Default())

	routes.UserRoutes(server)
	logger.Info("", "User Routes registered successfully")

	if err := server.Run(":8080"); err != nil {
		logger.Error("", "failed to start the server", err.Error())
	}
}
