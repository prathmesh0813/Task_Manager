package controller

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"task_manager/dao"
	"task_manager/middlewares"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Upload avatar
func UploadAvatar(c *gin.Context) {

	//check whether user is signin
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		utils.Logger.Warn("Unauthorized.User not authenticated", zap.Int64("userId", userId.(int64)))
		c.Set("response", nil)
		c.Set("message", "Unauthorized.User not authenticated")
		c.Set("error", true)
		c.Status(http.StatusUnauthorized)
		return
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		utils.Logger.Error("Invalid file", zap.Error(err), zap.Int64("userId", userId.(int64)))
		c.Set("response", nil)
		c.Set("message", "Invalid file")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	//validate the avatar
	fileExtension, err := utils.ValidateAvatar(fileHeader)
	if err != nil {
		utils.Logger.Error("Invalid file extension or file size", zap.Error(err), zap.Int64("userId", userId.(int64)))
		c.Set("response", nil)
		c.Set("message", err.Error())
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	//create filename to store in DB
	fileName := fmt.Sprintf("avatar_%s%s", strconv.FormatInt(userId.(int64), 10), fileExtension)

	file, err := fileHeader.Open()
	if err != nil {
		utils.Logger.Error("Failed to open uploded file", zap.Error(err), zap.Int64("userId", userId.(int64)))
		c.Set("response", nil)
		c.Set("message", "Failed to open uploded file")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		utils.Logger.Error("Failed to read uploded file", zap.Error(err), zap.Int64("userId", userId.(int64)))
		c.Set("response", nil)
		c.Set("message", "Failed to read uploded file")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	//checks whether the avatar is allready uploaded or not
	_, err = dao.ReadAvatar(userId.(int64))
	if err == nil {
		//If allready uploaded then it update avatar
		err = dao.UpdateAvatar(userId.(int64), content)
		if err != nil {
			//utils.StandardResponse(c, http.StatusInternalServerError, "Failed to save avatar", true, nil)
			utils.Logger.Error("Failed to save avatar", zap.Error(err), zap.Int64("userId", userId.(int64)))
			c.Set("response", nil)
			c.Set("message", "Failed to save avatar")
			c.Set("error", true)
			c.Status(http.StatusInternalServerError)
			return
		}

		utils.Logger.Info("Avatar updated", zap.Int64("userId", userId.(int64)))
		c.Set("response", nil)
		c.Set("message", "Avatar updated")
		c.Set("error", false)
		c.Status(http.StatusOK)
	} else {
		//If not uploaded then upload the avatar
		err = dao.SaveAvatar(userId.(int64), content, fileName)
		if err != nil {
			utils.Logger.Error("Failed to upload avatar", zap.Error(err), zap.Int64("userId", userId.(int64)))
			c.Set("response", nil)
			c.Set("message", "Failed to upload avatar")
			c.Set("error", true)
			c.Status(http.StatusInternalServerError)
			return
		}
	}
	defer file.Close()
	utils.Logger.Info("Avatar uploaded", zap.Int64("userId", userId.(int64)))
	c.Set("response", nil)
	c.Set("message", "Avatar uploaded")
	c.Set("error", false)
	c.Status(http.StatusOK)
}
