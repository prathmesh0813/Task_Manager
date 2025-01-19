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

//authentication middleware
func Authenticate(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		utils.Logger.Warn("Authorization token is missing", zap.String("method", c.Request.Method), zap.String("url", c.Request.URL.String()))
		c.JSON(http.StatusUnauthorized, gin.H{"message": "token not found", "error": true, "data": nil})
		c.Abort()
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")

	userId, err := utils.VerifyJwtToken(token)
	if err != nil {
		utils.Logger.Error("Failed to verify user token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized", "error": true, "data": nil})
		c.Abort()
		return
	}

	c.Set("token", token)
	c.Set("userId", userId)

	c.Next()
}

//check token in db
func CheckTokenPresent(c *gin.Context) error {
	token := c.Request.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	var dbToken dao.Token

	err := dao.DB.Where("user_token = ?", token).First(&dbToken).Error
	if err != nil {
		utils.Logger.Error("Session expired or token not found", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Session Expired: User has to log in", "error": true, "data": nil})
	}

	utils.Logger.Info("Token found in the database", zap.String("tokenId", fmt.Sprintf("%d", dbToken.ID)))
	return err
}
