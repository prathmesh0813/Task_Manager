package main

import (
	dao "task_manager/DAO"
	"task_manager/router"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
)

func main() {

	utils.InitLogger()
	dao.InitDB()
	server := gin.Default()
	router.RegisterRoutes(server)
	server.Run(":8080")
	// fmt.Println("Hello")
	// fmt.Println("Hello2")
}
