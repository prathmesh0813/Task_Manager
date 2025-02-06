package controller

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"task_manager/dao"
	"task_manager/logger"
	"task_manager/middlewares"
	"task_manager/utils"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// Upload avatar
func UploadAvatar(c *gin.Context) {
	requestID := requestid.Get(c)

	//check whether user is signin
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		logger.Warn(requestID, "Unauthorized, user not authenticated", "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "unauthorized, user not authenticated", true, http.StatusUnauthorized)
		return
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		logger.Error(requestID, "Invalid file", err.Error(), "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "invalid file", true, http.StatusBadRequest)
		return
	}

	//validate the avatar
	fileExtension, err := utils.ValidateAvatar(fileHeader)
	if err != nil {
		logger.Error(requestID, "Invalid file extension or file size", err.Error(), "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, err.Error(), true, http.StatusBadRequest)
		return
	}

	//create filename to store in DB
	fileName := fmt.Sprintf("avatar_%s%s", strconv.FormatInt(userId.(int64), 10), fileExtension)

	file, err := fileHeader.Open()
	if err != nil {
		logger.Error(requestID, "failed to open uploaded file", err.Error(), "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "failed to open uploaded file", true, http.StatusInternalServerError)
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		logger.Error(requestID, "failed to read uploaded file", err.Error(), "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "failed to read uploaded file", true, http.StatusInternalServerError)
		return
	}

	//checks whether the avatar is allready uploaded or not
	_, err = dao.ReadAvatar(userId.(int64))
	if err == nil {
		//If allready uploaded then it update avatar
		err = dao.UpdateAvatar(userId.(int64), content)
		if err != nil {
			logger.Error(requestID, "failed to save avatar", err.Error(), "userID: "+strconv.Itoa(int(userId.(int64))))
			utils.SetResponse(c, requestID, nil, "failed to save avatar", true, http.StatusInternalServerError)
			return
		}

		logger.Info(requestID, "Avatar updated", "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "avatar updated successfully", false, http.StatusOK)

	} else {
		//If not uploaded then upload the avatar
		err = dao.SaveAvatar(userId.(int64), content, fileName)
		if err != nil {
			logger.Error(requestID, "failed to upload avatar", err.Error(), "userID: "+strconv.Itoa(int(userId.(int64))))
			utils.SetResponse(c, requestID, nil, "failed to upload avatar", true, http.StatusInternalServerError)
			return
		}
	}
	defer file.Close()

	logger.Info(requestID, "Avatar uploaded successfully", "userID: "+strconv.Itoa(int(userId.(int64))))
	utils.SetResponse(c, requestID, nil, "avatar uploaded successfully", false, http.StatusOK)
}

// read avatar
func ReadAvatar(c *gin.Context) {
	requestID := requestid.Get(c)
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error(requestID, "failed to parse avatar Id", c.Param("id"), err.Error())
		utils.SetResponse(c, requestID, nil, "could not parse avatar id", true, http.StatusBadRequest)
		return
	}

	avatar, err := dao.ReadAvatar(userId)
	if err != nil {
		logger.Error(requestID, "failed to read avatar", err.Error())
		utils.SetResponse(c, requestID, nil, "failed to read avatar", true, http.StatusInternalServerError)
		return
	}

	logger.Info(requestID, "avatar fetched successfully")
	utils.SetResponse(c, requestID, string(avatar.Data), "Avatar fetch successfully", false, http.StatusOK)
}

// delete user avatar
func DeleteAvatar(c *gin.Context) {
	requestID := requestid.Get(c)
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		logger.Error(requestID, "User Id not found in context", "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "unauthorized, user not authenticated", true, http.StatusUnauthorized)
		return
	}

	_, err = dao.ReadAvatar(userId.(int64))
	if err != nil {
		logger.Error(requestID, "failed to read avatar", err.Error(), "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "no avatar present to delete", true, http.StatusInternalServerError)
		return
	}

	err = dao.DeleteAvatar(userId.(int64))
	if err != nil {
		logger.Error(requestID, "failed to delete avatar", err.Error(), "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "failed to delete user avatar", true, http.StatusInternalServerError)
		return
	}

	logger.Info(requestID, "avatar deleted successfully")
	utils.SetResponse(c, requestID, nil, "avatar deleted successfully", false, http.StatusOK)
}
