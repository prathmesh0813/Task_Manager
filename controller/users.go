package controller

import (
	"bytes"
	"encoding/json"
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

// user signup
func SignUp(c *gin.Context) {
	requestID := requestid.Get(c)
	var user models.User

	bodyBytes, _ := io.ReadAll(c.Request.Body)
	requestBody := string(bodyBytes)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var requestData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &requestData); err == nil {
		if _, ok := requestData["password"]; ok {
			requestData["password"] = "******"
		}

		maskedJSON, _ := json.Marshal(requestData)
		requestBody = string(maskedJSON)
	}

	err := c.ShouldBindJSON(&user)
	if err != nil {
		logger.Error(requestID, "cannot parsed the requested data", err.Error(), "userID: "+strconv.Itoa(int(user.ID)), requestBody)
		utils.SetResponse(c, requestID, nil, "cannot parsed the requested data", true, http.StatusBadRequest)
		return
	}

	//validate User credentials
	err = utils.ValidateDetails(user.Name, user.Email, user.Mobile_No, user.Gender, user.Password)
	if err != nil {
		logger.Error(requestID, "Unable to validate user details", err.Error(), "userID: "+strconv.Itoa(int(user.ID)), requestBody)
		utils.SetResponse(c, requestID, nil, err.Error(), true, http.StatusBadRequest)
		return
	}

	user.Gender = strings.ToLower(user.Gender)

	//Save tokens
	uid, userToken, refreshToken, err := dao.SaveUser(dao.DB, &user)
	if err != nil {
		logger.Error(requestID, "Unable to save user.User already exists.", err.Error(), "userID: "+strconv.Itoa(int(user.ID)), requestBody)
		utils.SetResponse(c, requestID, nil, "Unable to save user.User already exists.", true, http.StatusBadRequest)
		return
	}

	logger.Info(requestID, "User registered successfully", "userID: "+strconv.Itoa(int(uid)), requestBody)
	utils.SetResponse(c, requestID, gin.H{"refresh_token": refreshToken, "user_token": userToken}, "User registered successfully", false, http.StatusCreated)
}

// user sign in
func SignIn(c *gin.Context) {

	requestID := requestid.Get(c)
	var login models.Login

	bodyBytes, _ := io.ReadAll(c.Request.Body)
	requestBody := string(bodyBytes)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var requestData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &requestData); err == nil {
		if _, ok := requestData["password"]; ok {
			requestData["password"] = "******"
		}

		maskedJSON, _ := json.Marshal(requestData)
		requestBody = string(maskedJSON)
	}

	err := c.ShouldBindJSON(&login)
	if err != nil {
		logger.Warn(requestID, "failed to parse login request", err.Error(), "userID: "+strconv.Itoa(int(login.ID)), requestBody)
		utils.SetResponse(c, requestID, nil, "username and password required", true, http.StatusBadRequest)
		return
	}

	//validate credentials
	err = dao.ValidateCredentials(&login)
	if err != nil {
		logger.Warn(requestID, "Authentication failed", err.Error(), requestBody)
		utils.SetResponse(c, requestID, nil, "incorrect username or password", true, http.StatusBadRequest)
		return
	}

	//generating token
	userToken, refreshToken, err := utils.GenerateTokens(login.ID)
	if err != nil {
		logger.Error(requestID, "failed to generate tokens", "userID: "+strconv.Itoa(int(login.ID)), err.Error(), requestBody)
		utils.SetResponse(c, requestID, nil, "user login failed", true, http.StatusBadRequest)
		return
	}

	//save pair of token
	err = dao.SaveToken(login.ID, userToken, refreshToken)
	if err != nil {
		logger.Error(requestID, "failed to save tokens", "userID: "+strconv.Itoa(int(login.ID)), err.Error(), requestBody)
		utils.SetResponse(c, requestID, nil, "user login failed", true, http.StatusBadRequest)
		return
	}

	logger.Info(requestID, "user signed in successfully", "userID: "+strconv.Itoa(int(login.ID)), requestBody)
	utils.SetResponse(c, requestID, gin.H{"refresh_token": refreshToken, "user_token": userToken}, "user sign in successfully", false, http.StatusCreated)
}

// Fetch the user details
func GetUser(c *gin.Context) {
	requestID := requestid.Get(c)

	userId, exists := c.Get("userId")
	if !exists {
		logger.Warn(requestID, "Unauthorized, user not authenticated", "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "unauthorized, user not authenticated", true, http.StatusUnauthorized)
		return
	}

	//checks whether user is signin or not
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		logger.Warn(requestID, "session expired or token not found", "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "session expired or token not found", true, http.StatusBadRequest)
		return
	}

	//Get user by ID
	user, err := dao.GetUserById(userId.(int64))
	if err != nil {
		logger.Error(requestID, "could not fetch user", err.Error(), "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "could not fetch user", true, http.StatusBadRequest)
		return
	}

	user.Avatar = "/user/avatar/" + strconv.FormatInt(userId.(int64), 10)

	logger.Info(requestID, "User fetched successfully", "userID: "+strconv.Itoa(int(userId.(int64))))
	utils.SetResponse(c, requestID, user, "user fetched successfully", false, http.StatusOK)

}

func RefreshTokenHandler(c *gin.Context) {
	requestID := requestid.Get(c)
	// Get the refresh token from the request header
	err := middlewares.CheckRefreshToken(c)
	if err != nil {
		logger.Error(requestID, "refresh token required", "")
		utils.SetResponse(c, requestID, nil, "refresh token required", true, http.StatusUnauthorized)
		return
	}

	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		logger.Error(requestID, "refresh token required", "")
		utils.SetResponse(c, requestID, nil, "refresh token required", true, http.StatusUnauthorized)
		return
	}

	// Verify the refresh token
	userId, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil {
		err = dao.DeleteRefreshToken(refreshToken)
		if err != nil {
			logger.Error(requestID, "failed to delete refresh token", err.Error())
			utils.SetResponse(c, requestID, nil, "failed to delete refresh token", true, http.StatusBadRequest)
			return
		}

		logger.Error(requestID, "invalid refresh token", err.Error())
		utils.SetResponse(c, requestID, nil, "invalid refresh token", true, http.StatusUnauthorized)
		return
	}

	// Generate a new access token
	newUserToken, newRefreshToken, err := utils.GenerateTokens(userId)
	if err != nil {
		logger.Error(requestID, "failed to generate tokens", "userID: "+strconv.Itoa(int(userId)), err.Error())
		utils.SetResponse(c, requestID, nil, "user login failed", true, http.StatusBadRequest)
		return
	}

	err = dao.SaveToken(userId, newUserToken, newRefreshToken)
	if err != nil {
		logger.Error(requestID, "could not save token", "userID: "+strconv.Itoa(int(userId)), err.Error())
		utils.SetResponse(c, requestID, nil, "could not save tokens", true, http.StatusBadRequest)
		return
	}

	err = dao.DeleteRefreshToken(refreshToken)
	if err != nil {
		logger.Error(requestID, "failed to delete refresh token", "userID: "+strconv.Itoa(int(userId)), err.Error())
		utils.SetResponse(c, requestID, nil, "failed to delete refresh token", true, http.StatusBadRequest)
		return
	}

	// Return the new access token to the client
	logger.Info(requestID, "token refreshed successfully", "userID: "+strconv.Itoa(int(userId)))
	utils.SetResponse(c, requestID, gin.H{"refresh_token": newRefreshToken, "user_token": newUserToken}, "token refreshed successfully", false, http.StatusOK)
}

// Updates user details
func UpdateUser(c *gin.Context) {
	requestID := requestid.Get(c)
	var req models.UpdateUserRequest

	bodyBytes, _ := io.ReadAll(c.Request.Body)
	requestBody := string(bodyBytes)

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	userId, exists := c.Get("userId")
	if !exists {
		logger.Warn(requestID, "Unauthorized, user not authenticated", "userID: "+strconv.Itoa(int(userId.(int64))), requestBody)
		utils.SetResponse(c, requestID, nil, "unauthorized, user not authenticated", true, http.StatusUnauthorized)
		return
	}
	//checks whether user is signin or not
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		logger.Warn(requestID, "session expired or token not found", "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "session expired or token not found", true, http.StatusBadRequest)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(requestID, "Invalid request body", err.Error(), "userID: "+strconv.Itoa(int(userId.(int64))), requestBody)
		utils.SetResponse(c, requestID, nil, "invalid request body", true, http.StatusBadRequest)
		return
	}

	//Validates the user details taken from user response to update details
	err = utils.ValidateUser(req.Name, req.Mobile_No)
	if err != nil {
		logger.Error(requestID, "Unable to validate credentials", err.Error(), "userID: "+strconv.Itoa(int(userId.(int64))), requestBody)
		utils.SetResponse(c, requestID, nil, err.Error(), true, http.StatusBadRequest)
		return
	}

	//Updates user details
	err = dao.UpdateUserDetails(userId.(int64), req)
	if err != nil {
		logger.Error(requestID, "failed to update user", err.Error(), "userID: "+strconv.Itoa(int(userId.(int64))), requestBody)
		utils.SetResponse(c, requestID, nil, "failed to update user", true, http.StatusBadRequest)
		return
	}

	logger.Info(requestID, "User details updated successfully", "userID: "+strconv.Itoa(int(userId.(int64))), requestBody)
	utils.SetResponse(c, requestID, nil, "user details updated successfully", false, http.StatusOK)

}

// Update password
func UpdatePassword(c *gin.Context) {
	requestID := requestid.Get(c)
	var req models.UpdatePasswordRequest
	userIdFromToken, exists := c.Get("userId")
	if !exists {
		logger.Warn(requestID, "Unauthorized, user not authenticated", "userID: "+strconv.Itoa(int(userIdFromToken.(int64))))
		utils.SetResponse(c, requestID, nil, "unauthorized, user not authenticated", true, http.StatusUnauthorized)
		return
	}
	//checks whether user is signin or not
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		logger.Warn(requestID, "session expired or token not found", "userID: "+strconv.Itoa(int(userIdFromToken.(int64))))
		utils.SetResponse(c, requestID, nil, "session expired or token not found", true, http.StatusBadRequest)
		return
	}

	token := c.Request.Header.Get("Authorization")

	token = strings.TrimPrefix(token, "Bearer ")

	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(requestID, "Invalid request body", err.Error(), "userID: "+strconv.Itoa(int(userIdFromToken.(int64))))
		utils.SetResponse(c, requestID, nil, "invalid request body", true, http.StatusBadRequest)
		return
	}

	//Validate whether enter password is in correct format or not
	err = utils.ValidatePassword(req.NewPassword)
	if err != nil {
		logger.Error(requestID, "unable to validate credentials", err.Error(), "userID: "+strconv.Itoa(int(userIdFromToken.(int64))))
		utils.SetResponse(c, requestID, nil, err.Error(), true, http.StatusBadRequest)
		return
	}

	//Check whether user is present or not in db to chng password
	user, err := dao.GetUserByIdPassChng(userIdFromToken.(int64))
	if err != nil {
		logger.Error(requestID, "User not found", err.Error(), "userID: "+strconv.Itoa(int(userIdFromToken.(int64))))
		utils.SetResponse(c, requestID, nil, "user not found", true, http.StatusBadRequest)
		return
	}

	//checks the user enter old password is correct with password hash which is present in DB
	passwordIsValid := utils.CheckPasswordHash(req.OldPassword, user.Password)
	if !passwordIsValid {
		logger.Error(requestID, "incorrect old password", "userID: "+strconv.Itoa(int(userIdFromToken.(int64))))
		utils.SetResponse(c, requestID, nil, "incorrect old password", true, http.StatusBadRequest)
		return
	}

	//Hash the new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		logger.Error(requestID, "failed to hashed password", err.Error(), "userID: "+strconv.Itoa(int(userIdFromToken.(int64))))
		utils.SetResponse(c, requestID, nil, "failed to hashed password", true, http.StatusBadRequest)
		return
	}

	//Update password by userId
	err = dao.UpdatePassById(userIdFromToken.(int64), hashedPassword)
	if err != nil {
		logger.Error(requestID, "failed to update password", err.Error(), "userID: "+strconv.Itoa(int(userIdFromToken.(int64))))
		utils.SetResponse(c, requestID, nil, "failed to update password", true, http.StatusBadRequest)
		return
	}

	//Signout all other user except the user which updates the password
	err = dao.DeleteTokenById(userIdFromToken.(int64), token)
	if err != nil {
		logger.Error(requestID, "failed to delete token", err.Error(), "userID: "+strconv.Itoa(int(userIdFromToken.(int64))))
		utils.SetResponse(c, requestID, nil, "failed to delete token", true, http.StatusBadRequest)
		return
	}

	logger.Info(requestID, "Password updated successfully", "userID: "+strconv.Itoa(int(userIdFromToken.(int64))))
	utils.SetResponse(c, requestID, nil, "password updated successfully", false, http.StatusOK)

}

// Signout user
func SignOut(c *gin.Context) {
	requestID := requestid.Get(c)

	//checks whether user is signin or not
	err := middlewares.CheckTokenPresent(c)
	userId, exists := c.Get("userId")
	if !exists {
		logger.Warn(requestID, "Unauthorized, user not authenticated", "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "unauthorized, user not authenticated", true, http.StatusUnauthorized)
		return
	}
	if err != nil {
		logger.Warn(requestID, "session expired or token not found", "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "session expired or token not found", true, http.StatusBadRequest)
		return
	}

	tokenString := strings.TrimSpace(c.GetHeader("Authorization"))
	if tokenString == "" {
		logger.Error(requestID, "Unable to validate user details", "")
		utils.SetResponse(c, requestID, nil, "token not provided", true, http.StatusUnauthorized)
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	//takes query parameters
	allParam := c.DefaultQuery("all", "false")

	all, err := strconv.ParseBool(allParam)
	if err != nil {
		logger.Error(requestID, "Invalid query parameter for 'all'.It must be true or false", err.Error(), "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "invalid query parameter for 'all'. It must be true or false", true, http.StatusBadRequest)
		return

	}

	//if all query param value is true

	if all {
		//signout all users
		err = dao.SignOutAllUsers(userId.(int64))
		if err != nil {
			logger.Error(requestID, "failed to signout from all devices", err.Error(), "userID: "+strconv.Itoa(int(userId.(int64))))
			utils.SetResponse(c, requestID, nil, "failed to signout from all devices", true, http.StatusBadRequest)
			return
		}

		logger.Info(requestID, "Signout from all devices successfully", "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "signout from all devices successfully", false, http.StatusOK)

	} else {
		//Sign out single user if all query param value is false
		err = dao.DeleteToken(tokenString)
		if err != nil {
			logger.Error(requestID, "failed to signout", err.Error(), "userID: "+strconv.Itoa(int(userId.(int64))))
			utils.SetResponse(c, requestID, nil, "failed to sign out", true, http.StatusBadRequest)
			return
		}

		logger.Info(requestID, "User sign out successfully", "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "user sign out successfully", false, http.StatusOK)
	}
}
