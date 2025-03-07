package controller

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"
	"task_manager/dao"
	"task_manager/logger"
	"task_manager/middlewares"
	"task_manager/models"
	"task_manager/utils"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func CreateChatHead(c *gin.Context) {
	var email models.Emails

	requestID := requestid.Get(c)

	userId, exists := c.Get("userId")
	if !exists {
		logger.Warn(requestID, "Unauthorized, user not authorized", "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "Unauthorized, user not authorized", true, http.StatusUnauthorized)
		return
	}

	//checks whether user is signin or not
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		logger.Warn(requestID, "session expired or token not found", "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "session expired or token not found", true, http.StatusBadRequest)
		return
	}

	bodyBytes, _ := io.ReadAll(c.Request.Body)
	requestBody := string(bodyBytes)

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	err = c.ShouldBindJSON(&email)
	if err != nil {
		logger.Error(requestID, "cannot parsed the requested data", "userID: "+strconv.Itoa(int(userId.(int64))), requestBody)
		utils.SetResponse(c, requestID, nil, "cannot parsed the requested data", true, http.StatusBadRequest)
		return
	}

	isValidEmail, emailID := utils.EmailValidation(email.Username)

	invalidEmailStr := strings.Join(emailID, ", ")
	if !isValidEmail {
		logger.Error(requestID, "Email is in invalid format", "userID: "+strconv.Itoa(int(userId.(int64))), requestBody)
		utils.SetResponse(c, requestID, nil, "Email is in invalid format "+invalidEmailStr, true, http.StatusBadRequest)
		return
	}

	err, emailId := dao.CheckEmailPresent(email.Username)
	misingEmail := strings.Join(emailId, ", ")
	if err != nil {
		logger.Error(requestID, "User does not exist", "userID: "+strconv.Itoa(int(userId.(int64))), requestBody)
		utils.SetResponse(c, requestID, nil, "User does not exist "+misingEmail, true, http.StatusBadRequest)
		return
	}

	trimStr := strings.TrimSpace(email.GroupName)
	groupNameStr := strings.ReplaceAll(trimStr, " ", "")
	chatHead := strings.ToLower(groupNameStr)

	_, err = utils.CreateHash(chatHead)
	if err != nil {
		logger.Error(requestID, "Chat head not hashed", "userID: "+strconv.Itoa(int(userId.(int64))), requestBody)
		utils.SetResponse(c, requestID, nil, "Chat head not hashed", true, http.StatusBadRequest)
		return
	}

	logger.Info(requestID, "Email validate successfully", "userID: "+strconv.Itoa(int(userId.(int64))))
	utils.SetResponse(c, requestID, nil, "Email validate successfully", false, http.StatusOK)

}
