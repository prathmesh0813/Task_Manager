package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"task_manager/dao"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Set the response in gin context
func ResponseFormatter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		response, _ := c.Get("response")
		message, _ := c.Get("message")

		errorStatus, _ := c.Get("error")

		formattedResponse := gin.H{
			"message": message,
			"error":   errorStatus,
			"data":    response,
		}

		c.JSON(c.Writer.Status(), formattedResponse)

	}
}

// Checks whether the user is authenticated to perform the action
func Authenticate(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		utils.Logger.Warn("Authorization token is missing", zap.String("method", c.Request.Method), zap.String("url", c.Request.URL.String()))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized", "error": true, "data": nil})
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")

	userId, err := utils.VerifyJwtToken(token)
	if err != nil {
		utils.Logger.Error("Failed to verify user token", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized", "error": true, "data": nil})
		return
	}
	c.Set("token", token)
	c.Set("userId", userId)

	utils.Logger.Info("User authentication successfully", zap.String("userId", fmt.Sprintf("%d", userId)))
	c.Next()
}

// checks whether user is signin or not
func CheckTokenPresent(c *gin.Context) error {
	token := c.Request.Header.Get("Authorization")

	token = strings.TrimPrefix(token, "Bearer ")

	var dbToken dao.Token

	err := dao.DB.Where("user_token = ? ", token).First(&dbToken).Error
	if err != nil {
		utils.Logger.Warn("Session expired or token not found", zap.Error(err))
		c.Set("response", nil)
		c.Set("message", "Session Expired.User has to log in")
		c.Set("error", true)
		c.Status(http.StatusNotFound)

	}

	utils.Logger.Info("Token found in database", zap.String("tokenID", fmt.Sprintf("%d", dbToken.ID)))
	return err
}
