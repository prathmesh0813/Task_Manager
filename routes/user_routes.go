package routes

import (
	"task_manager/controller"
	"task_manager/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(server *gin.Engine) {
	route := server.Group("/user")

	route.POST("/signUp,", controller.SignUp, middlewares.ResponseFormatter())
	route.POST("/signIn", controller.SignIn, middlewares.ResponseFormatter())
	
}
