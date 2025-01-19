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

	route.GET("", middlewares.Authenticate, controller.GetUser, middlewares.ResponseFormatter())
	route.POST("/avatar", middlewares.Authenticate, controller.UploadAvatar, middlewares.ResponseFormatter())
	route.GET("/avatar/:id", controller.ReadAvatar, middlewares.ResponseFormatter())
	route.DELETE("/avatar", middlewares.Authenticate, controller.DeleteAvatar, middlewares.ResponseFormatter())
	route.POST("/refresh", controller.RefreshTokenHandler, middlewares.ResponseFormatter())
	route.PUT("/updateuser", middlewares.Authenticate, controller.UpdateUser, middlewares.ResponseFormatter())
	route.PUT("/updatepassword", middlewares.Authenticate, controller.UpdatePassword, middlewares.ResponseFormatter())
}
