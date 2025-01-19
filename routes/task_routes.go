package routes

import (
	"task_manager/controller"
	"task_manager/middlewares"

	"github.com/gin-gonic/gin"
)

func TaskRoutes(server *gin.Engine) {
	route := server.Group("/")

	route.POST("/tasks", middlewares.Authenticate, controller.CreateTask, middlewares.ResponseFormatter())
	route.GET("/tasks/:id", middlewares.Authenticate, controller.GetTask, middlewares.ResponseFormatter())
	route.GET("/tasks", middlewares.Authenticate, controller.GetTasksByQuery, middlewares.ResponseFormatter())
}
