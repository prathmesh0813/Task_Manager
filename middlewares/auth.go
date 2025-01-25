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

// check whether refresh token is present in db or not
func CheckRefreshToken(context *gin.Context) error {
	refreshToken := context.Request.Header.Get("Refresh-Token")
	if refreshToken == "" {
		utils.Logger.Error("Refresh token is missing in request header")
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh token required", "error": true, "data": nil})
		return fmt.Errorf("refresh token missing")
	}

	var dbToken dao.Token

	err := dao.DB.Where("refresh_token = ?", refreshToken).First(&dbToken).Error
	if err != nil {
		utils.Logger.Error("Session expired or refresh token not found", zap.Error(err))
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired refresh token", "error": true, "data": nil})
		return err
	}

	utils.Logger.Info("Token found in the database", zap.String("tokenId", fmt.Sprintf("%d", dbToken.ID)))
	return err
}
