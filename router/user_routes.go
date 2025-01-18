package router

import (
	"task_manager/controller"
	"task_manager/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {

	userGroup := server.Group("/user")
	userGroup.POST("/signUp", controller.SignUp, middlewares.ResponseFormatter())

}
