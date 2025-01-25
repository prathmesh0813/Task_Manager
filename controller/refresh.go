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
		utils.SetResponse(c, nil, "refresh token required", true, http.StatusUnauthorized)
		return
	}

	// Verify the refresh token
	userId, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil {
		err = dao.DeleteRefreshToken(refreshToken)
		if err != nil {
			utils.Logger.Error("failed to delete refresh token", zap.Error(err))
			utils.SetResponse(c, nil, "failed to delete refresh token", true, http.StatusNotFound)
			return
		}

		utils.Logger.Error("invalid refresh token", zap.Error(err))
		utils.SetResponse(c, nil, "invalid refresh token", true, http.StatusUnauthorized)
		return
	}

	newUserToken, newRefreshToken, err := utils.GenerateTokens(userId)
	if err != nil {
		utils.Logger.Error("error generating new access token", zap.Error(err))
		utils.SetResponse(c, nil, "error generating new access token", true, http.StatusInternalServerError)
		return
	}

	err = dao.SaveToken(userId, newUserToken, newRefreshToken)
	if err != nil {
		utils.Logger.Error("could not save token", zap.Error(err))
		utils.SetResponse(c, nil, "could not save token", true, http.StatusInternalServerError)
		return
	}

	err = dao.DeleteRefreshToken(refreshToken)
	if err != nil {
		utils.Logger.Error("failed to delete refresh token", zap.Error(err))
		utils.SetResponse(c, nil, "failed to delete refresh token", true, http.StatusNotFound)
		return
	}

	// Return the new access token to the client
	utils.Logger.Info("token refreshed successfully")
	utils.SetResponse(c, gin.H{"refresh_token": newRefreshToken, "user_token": newUserToken}, "token refreshed successfully", false, http.StatusOK)
}
