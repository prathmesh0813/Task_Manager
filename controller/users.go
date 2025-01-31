package controller

import (
	"net/http"
	"strconv"
	"strings"

	"task_manager/dao"
	"task_manager/middlewares"
	"task_manager/models"
	"task_manager/utils"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// User Signup
func SignUp(c *gin.Context) {
	requestId := requestid.Get(c)
	var user models.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		utils.Logger.Error("cannot parsed the requested data", zap.Error(err), zap.Int64("userId", user.ID))
		utils.SetResponse(c, requestId, nil, "failed to register the user", true, http.StatusBadRequest)
		return
	}

	//validate User credentials
	err = utils.ValidateDetails(user.Name, user.Email, user.Mobile_No, user.Gender, user.Password, user.Confirm_Password)
	if err != nil {
		utils.Logger.Error("User details not validate", zap.Error(err), zap.Int64("userId", user.ID))
		utils.SetResponse(c, requestId, nil, err.Error(), true, http.StatusBadRequest)
		return
	}

	user.Gender = strings.ToLower(user.Gender)

	//Save user in DB
	uid, err := dao.SaveUser(&user)
	if err != nil {
		utils.Logger.Error("Unable to save user. User already exist", zap.Error(err), zap.Int64("userId", user.ID))
		utils.SetResponse(c, requestId, nil, "user already exist", true, http.StatusBadRequest)
		return

	}

	//Generate pair of tokens
	userToken, refreshToken, err := utils.GenerateTokens(uid)
	if err != nil {
		utils.Logger.Error("error generating pair of tokens", zap.Error(err))
		utils.SetResponse(c, requestId, nil, "error generating pair of tokens", true, http.StatusInternalServerError)
		return
	}

	//Save tokens
	err = dao.SaveToken(uid, userToken, refreshToken)
	if err != nil {
		utils.Logger.Error("failed to save the token", zap.Error(err), zap.Int64("userId", user.ID))
		utils.SetResponse(c, requestId, nil, "failed to register the user", true, http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("User registered successfully", zap.Int64("userId", uid))
	utils.SetResponse(c, requestId, gin.H{"refresh_token": refreshToken, "user_token": userToken}, "User registered successfully", false, http.StatusCreated)

}

// user sign in
func SignIn(c *gin.Context) {
	requestId := requestid.Get(c)
	var login models.Login

	err := c.ShouldBindJSON(&login)
	if err != nil {
		utils.Logger.Warn("Failed to parse login request", zap.Error(err))
		utils.SetResponse(c, requestId, nil, "username and password required", true, http.StatusBadRequest)
		return
	}

	//validate credentials
	err = dao.ValidateCredentials(&login)
	if err != nil {
		utils.Logger.Warn("Authentication failed", zap.Error(err))
		utils.SetResponse(c, requestId, nil, "incorrect username or password", true, http.StatusBadRequest)
		return
	}
	//Generate pair of tokens
	userToken, refreshToken, err := utils.GenerateTokens(login.ID)
	if err != nil {
		utils.Logger.Error("error generating pair of tokens", zap.Error(err))
		utils.SetResponse(c, requestId, nil, "error generating pair of tokens", true, http.StatusInternalServerError)
		return
	}

	//save pair of token
	err = dao.SaveToken(login.ID, userToken, refreshToken)
	if err != nil {
		utils.Logger.Error("Failed to save token", zap.Int64("userId", login.ID), zap.Error(err))
		utils.SetResponse(c, requestId, nil, "user login failed", true, http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("User signed in successfully", zap.Int64("userId", login.ID))
	utils.SetResponse(c, requestId, gin.H{"refresh_token": refreshToken, "user_token": userToken}, "user sign in successfully", false, http.StatusCreated)

}

// Fetch the user details
func GetUser(c *gin.Context) {
	requestId := requestid.Get(c)

	//checks whether user is signin or not
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		utils.Logger.Warn("Unauthorized.User not authenticated", zap.Error(err), zap.Int64("userId", userId.(int64)))
		utils.SetResponse(c, requestId, nil, "Unauthorized.User not authenticated", true, http.StatusUnauthorized)
		return
	}

	//Get user by ID
	user, err := dao.GetUserById(userId.(int64))
	if err != nil {
		utils.Logger.Error("Could not fetch user", zap.Error(err), zap.Int64("userId", userId.(int64)))
		utils.SetResponse(c, requestId, nil, "Could not fetch user", true, http.StatusInternalServerError)
		return
	}

	user.Avatar = "/user/avatar/" + strconv.FormatInt(userId.(int64), 10)
	utils.Logger.Info("User Fetch successfully", zap.Error(err), zap.Int64("userId", userId.(int64)))
	utils.SetResponse(c, requestId, gin.H{"user": user}, "User Fetch successfully", false, http.StatusOK)
}

// Updates user details
func UpdateUser(c *gin.Context) {
	requestId := requestid.Get(c)
	var req models.UpdateUserRequest
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}
	userId, exists := c.Get("userId")
	if !exists {
		utils.Logger.Warn("Unauthorized .User not authenticated", zap.Error(err), zap.Int64("userId", userId.(int64)))
		utils.SetResponse(c, requestId, nil, "Unauthorized .User not authenticated", true, http.StatusUnauthorized)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Error("Invalid request body", zap.Error(err), zap.Int64("userId", userId.(int64)))
		utils.SetResponse(c, requestId, nil, "Invalid request body", true, http.StatusBadRequest)
		return
	}

	//Validates the user details taken from user response to update details
	err = utils.ValidateUser(req.Name, req.Mobile_No)
	if err != nil {
		utils.Logger.Error("Credentials are not validate", zap.Error(err), zap.Int64("userId", userId.(int64)))
		utils.SetResponse(c, requestId, nil, err.Error(), true, http.StatusBadRequest)
		return
	}

	//Updates user details
	err = dao.UpdateUserDetails(userId.(int64), req)
	if err != nil {
		utils.Logger.Error("Failed to update user", zap.Error(err), zap.Int64("userId", userId.(int64)))
		utils.SetResponse(c, requestId, nil, "Failed to update user", true, http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("User details updated successfully", zap.Int64("userId", userId.(int64)))
	utils.SetResponse(c, requestId, nil, "User details updated successfully", false, http.StatusOK)

}

// Update password
func UpdatePassword(c *gin.Context) {
	requestId := requestid.Get(c)
	var req models.UpdatePasswordRequest
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	token := c.Request.Header.Get("Authorization")

	token = strings.TrimPrefix(token, "Bearer ")

	userIdFromToken, exists := c.Get("userId")
	if !exists {
		utils.Logger.Warn("Unauthorized.User not authenticated", zap.Error(err), zap.Int64("userId", userIdFromToken.(int64)))
		utils.SetResponse(c, requestId, nil, "Unauthorized.User not authenticated", true, http.StatusUnauthorized)
		return
	}
	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Error("Invalid request body", zap.Error(err), zap.Int64("userId", userIdFromToken.(int64)))
		utils.SetResponse(c, requestId, nil, "Invalid request body", true, http.StatusBadRequest)
		return
	}

	//Validate whether enter password is in correct format or not
	err = utils.ValidatePassword(req.NewPassword)
	if err != nil {
		utils.Logger.Error("Credentials are not validate", zap.Error(err), zap.Int64("userId", userIdFromToken.(int64)))
		utils.SetResponse(c, requestId, nil, err.Error(), true, http.StatusBadRequest)
		return
	}

	//Check whether user is present or not in db to chng password
	user, err := dao.GetUserByIdPassChng(userIdFromToken.(int64))
	if err != nil {
		utils.Logger.Error("User not found", zap.Error(err), zap.Int64("userId", userIdFromToken.(int64)))
		utils.SetResponse(c, requestId, nil, "User not found", true, http.StatusNotFound)
		return
	}

	//checks the user enter old password is correct with password hash which is present in DB
	passwordIsValid := utils.CheckPasswordHash(req.OldPassword, user.Password)
	if !passwordIsValid {
		utils.Logger.Error("Incorrect old password", zap.Error(err), zap.Int64("userId", userIdFromToken.(int64)))
		utils.SetResponse(c, requestId, nil, "Incorrect old password", true, http.StatusBadRequest)
		return
	}

	//Hash the new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.Logger.Error("Failed to hashed password", zap.Error(err), zap.Int64("userId", userIdFromToken.(int64)))
		utils.SetResponse(c, requestId, nil, "Failed to hashed password", true, http.StatusInternalServerError)
		return
	}

	//Update password by userId
	err = dao.UpdatePassById(userIdFromToken.(int64), hashedPassword)
	if err != nil {
		utils.Logger.Error("Failed to update password", zap.Error(err), zap.Int64("userId", userIdFromToken.(int64)))
		utils.SetResponse(c, requestId, nil, "Failed to update password", true, http.StatusInternalServerError)
		return
	}

	//Signout all other user except the user which updates the password
	err = dao.DeleteTokenById(userIdFromToken.(int64), token)
	if err != nil {
		utils.Logger.Error("Failed to delete token", zap.Error(err), zap.Int64("userId", userIdFromToken.(int64)))
		utils.SetResponse(c, requestId, nil, "Failed to delete token", true, http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("Password Updated Successfully.", zap.Int64("userId", userIdFromToken.(int64)))
	utils.SetResponse(c, requestId, nil, "Password Updated Successfully.", false, http.StatusOK)
}

func RefreshTokenHandler(c *gin.Context) {
	requestId := requestid.Get(c)
	// Get the refresh token from the request header
	err := middlewares.CheckRefreshToken(c)
	if err != nil {
		return
	}

	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		utils.Logger.Error("refresh token required", zap.Error(err))
		utils.SetResponse(c, requestId, nil, "refresh token required", true, http.StatusUnauthorized)
		return
	}

	// Verify the refresh token
	userId, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil {
		err = dao.DeleteRefreshToken(refreshToken)
		if err != nil {
			utils.Logger.Error("failed to delete refresh token", zap.Error(err))
			utils.SetResponse(c, requestId, nil, "failed to delete refresh token", true, http.StatusNotFound)
			return
		}

		utils.Logger.Error("invalid refresh token", zap.Error(err))
		utils.SetResponse(c, requestId, nil, "invalid refresh token", true, http.StatusUnauthorized)
		return
	}

	newUserToken, newRefreshToken, err := utils.GenerateTokens(userId)
	if err != nil {
		utils.Logger.Error("error generating new access token", zap.Error(err))
		utils.SetResponse(c, requestId, nil, "error generating new access token", true, http.StatusInternalServerError)
		return
	}

	err = dao.SaveToken(userId, newUserToken, newRefreshToken)
	if err != nil {
		utils.Logger.Error("could not save token", zap.Error(err))
		utils.SetResponse(c, requestId, nil, "could not save token", true, http.StatusInternalServerError)
		return
	}

	err = dao.DeleteRefreshToken(refreshToken)
	if err != nil {
		utils.Logger.Error("failed to delete refresh token", zap.Error(err))
		utils.SetResponse(c, requestId, nil, "failed to delete refresh token", true, http.StatusNotFound)
		return
	}

	// Return the new access token to the client
	utils.Logger.Info("token refreshed successfully")
	utils.SetResponse(c, requestId, gin.H{"refresh_token": newRefreshToken, "user_token": newUserToken}, "token refreshed successfully", false, http.StatusOK)
}

// Signout user
func SignOut(c *gin.Context) {
	requestId := requestid.Get(c)

	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	tokenString := strings.TrimSpace(c.GetHeader("Authorization"))
	if tokenString == "" {
		utils.Logger.Error("Token not provided")
		utils.SetResponse(c, requestId, nil, "Token not provided", true, http.StatusUnauthorized)
		return
	}
	userId, exists := c.Get("userId")
	if !exists {
		utils.Logger.Warn("Unauthorized.User not authenticated", zap.Int64("userId", userId.(int64)))
		utils.SetResponse(c, requestId, nil, "Unauthorized.User not authenticated", true, http.StatusUnauthorized)
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	//takes query parameters
	allParam := c.DefaultQuery("all", "false")
	all, err := strconv.ParseBool(allParam)

	if err != nil {
		utils.Logger.Error("Invalid query parameter for 'all'. It must be true or false", zap.Error(err), zap.Int64("userId", userId.(int64)))
		utils.SetResponse(c, requestId, nil, "Invalid query parameter for 'all'. It must be true or false", true, http.StatusBadRequest)
		return

	}

	//if all query param value is true

	if all {
		//signout all users
		err = dao.SignOutAllUsers(userId.(int64))
		if err != nil {
			utils.Logger.Error("Failed to signout from all devices", zap.Error(err), zap.Int64("userId", userId.(int64)))
			utils.SetResponse(c, requestId, nil, "Failed to signout from all devices", true, http.StatusBadRequest)
		}
		utils.Logger.Info("Signout from all devices is successfully", zap.Int64("userId", userId.(int64)))
		utils.SetResponse(c, requestId, nil, "Signout from all devices is successfully", false, http.StatusOK)
	} else {
		//Sign out single user if all query param value is false
		err = dao.DeleteToken(tokenString)
		if err != nil {
			utils.Logger.Error("Failed to signout", zap.Error(err), zap.Int64("userId", userId.(int64)))
			utils.SetResponse(c, requestId, nil, "Failed to signout", true, http.StatusBadRequest)
			return
		}

		utils.Logger.Info("User sign out successfully", zap.Int64("userId", userId.(int64)))
		utils.SetResponse(c, requestId, nil, "User sign out successfully", false, http.StatusOK)
	}
}
