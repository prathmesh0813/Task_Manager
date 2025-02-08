package middlewares

import (
	"net/http"
	"strconv"
	"strings"
	"task_manager/dao"
	"task_manager/logger"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
)

// Checks whether the user is authenticated to perform the action
func Authenticate(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		logger.Warn("authorization-request-id", "Authorization token is missing", c.Request.Method, c.Request.URL.String())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization token is missing", "error": true, "data": nil})
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")

	userId, err := utils.VerifyJwtToken(token)
	if err != nil {
		logger.Error("", "failed to verify user token", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized", "error": true, "data": nil})
		return
	}

	c.Set("token", token)
	c.Set("userId", userId)

	logger.Info("requestID", "user authenticated successfully", strconv.Itoa(int(userId)))
	c.Next()
}

// checks whether user is signin or not
func CheckTokenPresent(c *gin.Context) error {
	token := c.Request.Header.Get("Authorization")

	token = strings.TrimPrefix(token, "Bearer ")

	var dbToken dao.Token

	err := dao.DB.Where("user_token = ? ", token).First(&dbToken).Error
	if err != nil {
		return err
	}

	logger.Info("requestID", "token found in the database", strconv.Itoa(int(dbToken.ID)))
	return err
}

// check whether refresh token is present in db or not
func CheckRefreshToken(context *gin.Context) error {
	refreshToken := context.Request.Header.Get("Refresh-Token")
	if refreshToken == "" {
		logger.Error("requestID", "Refresh token missing", "error")
		return nil

	}

	var dbToken dao.Token

	err := dao.DB.Where("refresh_token = ?", refreshToken).First(&dbToken).Error
	if err != nil {
		logger.Error("requestID", "session expired or refresh token not found", err.Error())
		return err
	}

	logger.Info("requestID", "refresh token found in the database", strconv.Itoa(int(dbToken.ID)))
	return nil
}
