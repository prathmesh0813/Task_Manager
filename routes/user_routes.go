package routes

import (
	"task_manager/controller"
	"task_manager/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(server *gin.Engine) {
	route := server.Group("/user")

	route.POST("/signup", controller.SignUp, middlewares.ResponseFormatter())
	route.POST("/signin", controller.SignIn, middlewares.ResponseFormatter())

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate, middlewares.ResponseFormatter())
	authenticated.GET("/user", controller.GetUser)

}
