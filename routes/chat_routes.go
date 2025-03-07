package routes

import (
	"task_manager/controller"
	"task_manager/middlewares"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(server *gin.Engine){
	route := server.Group("/chats",middlewares.RequestID())

	route.POST("/createchatheads", middlewares.Authenticate,controller.CreateChatHead,middlewares.ResponseFormatter())
}