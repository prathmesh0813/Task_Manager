package controller

import (
	"net/http"
	"task_manager/dao"
	"task_manager/middlewares"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RefreshTokenHandler(c *gin.Context) {
	// Get the refresh token from the request header
	err := middlewares.CheckRefreshToken(c)
	if err != nil {
		return
	}

	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		utils.Logger.Error("refresh token required", zap.Error(err))
		c.Set("response", nil)
		c.Set("message", "refresh token required")
		c.Set("error", true)
		c.Status(http.StatusUnauthorized)
		return
	}

	// Verify the refresh token
	userId, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil {
		err = dao.DeleteRefreshToken(refreshToken)
		if err != nil {
			utils.Logger.Error("failed to delete refresh token", zap.Error(err))

			c.Set("response", nil)
			c.Set("message", "failed to delete refresh token")
			c.Set("error", true)
			c.Status(http.StatusNotFound)
			return
		}

		utils.Logger.Error("invalid refresh token", zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "invalid refresh token")
		c.Set("error", true)
		c.Status(http.StatusUnauthorized)
		return
	}

	// Generate a new access token
	newUserToken, err := utils.GenerateJwtToken(userId)
	if err != nil {
		utils.Logger.Error("error generating new access token", zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "error generating new access token")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken(userId)
	if err != nil {
		utils.Logger.Error("error generating new refresh token", zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "error generating new refresh token")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	err = dao.SaveToken(userId, newUserToken, newRefreshToken)
	if err != nil {
		utils.Logger.Error("could not save token", zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "could not save token")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	err = dao.DeleteRefreshToken(refreshToken)
	if err != nil {
		utils.Logger.Error("failed to delete refresh token", zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "failed to delete refresh token")
		c.Set("error", true)
		c.Status(http.StatusNotFound)
		return
	}

	// Return the new access token to the client
	utils.Logger.Info("token refreshed successfully")

	c.Set("response",  gin.H{"refresh_token": newRefreshToken, "user_token": newUserToken})
	c.Set("message", "token refreshed successfully")
	c.Set("error", false)
	c.Status(http.StatusOK)
}
