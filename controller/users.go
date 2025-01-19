package controller

import (
	"net/http"
	"strconv"

	"task_manager/dao"
	"task_manager/middlewares"
	"task_manager/models"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// User Signup
func SignUp(c *gin.Context) {
	var user models.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		utils.Logger.Error("cannot parsed the requested data", zap.Error(err), zap.Int64("userId", user.ID))
		c.Set("response", nil)
		c.Set("message", "failed to register the user")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	//validate User credentials
	err = utils.ValidateDetails(user.Name, user.Email, user.Mobile_No, user.Gender, user.Password, user.Confirm_Password)
	if err != nil {
		utils.Logger.Error("User details no validate", zap.Error(err), zap.Int64("userId", user.ID))
		c.Set("response", nil)
		c.Set("message", err.Error())
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	//Save user in DB
	uid, err := dao.SaveUser(&user)
	if err != nil {
		utils.Logger.Error("Unable to save user. User already exist", zap.Error(err), zap.Int64("userId", user.ID))
		c.Set("response", nil)
		c.Set("message", "user already exist")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return

	}

	//Generate user token
	userToken, err := utils.GenerateJwtToken(uid)
	if err != nil {
		utils.Logger.Error("unable to generate the user token", zap.Error(err), zap.Int64("userId", user.ID))
		c.Set("response", nil)
		c.Set("message", "failed to register the user")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	//Generate Refresh token
	refreshToken, err := utils.GenerateRefreshToken(uid)
	if err != nil {
		utils.Logger.Error("unable to generate the refresh token", zap.Error(err), zap.Int64("userId", user.ID))
		c.Set("response", nil)
		c.Set("message", "failed to register the user")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	//Save tokens
	err = dao.SaveToken(uid, userToken, refreshToken)
	if err != nil {
		utils.Logger.Error("failed to save the token", zap.Error(err), zap.Int64("userId", user.ID))
		c.Set("response", nil)
		c.Set("message", "failed to register the user")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("User registered successfully", zap.Int64("userId", uid))
	c.Set("response", gin.H{"refresh_token": refreshToken, "user_token": userToken})
	c.Set("message", "User registered successfully")
	c.Set("error", false)
	c.Status(http.StatusCreated)
}

// user sign in
func SignIn(c *gin.Context) {
	var login models.Login

	err := c.ShouldBindJSON(&login)
	if err != nil {
		utils.Logger.Warn("Failed to parse login request", zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "username and password required")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	//validate credentials
	err = dao.ValidateCredentials(&login)
	if err != nil {
		utils.Logger.Warn("Authentication failed", zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "incorrect username or password")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	//generate user token
	userToken, err := utils.GenerateJwtToken(login.ID)
	if err != nil {
		utils.Logger.Error("Failed to generate user token", zap.Int64("userId", login.ID), zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "user login failed")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	//generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(login.ID)
	if err != nil {
		utils.Logger.Error("Failed to generate refresh token", zap.Int64("userId", login.ID), zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "user login failed")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	//save pair of token
	err = dao.SaveToken(login.ID, userToken, refreshToken)
	if err != nil {
		utils.Logger.Error("Failed to save token", zap.Int64("userId", login.ID), zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "user login failed")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("User signed in successfully", zap.Int64("userId", login.ID))

	c.Set("response", gin.H{"refresh_token": refreshToken, "user_token": userToken})
	c.Set("message", "user sign in successfully")
	c.Set("error", false)
	c.Status(http.StatusCreated)
}

// Fetch the user details
func GetUser(c *gin.Context) {

	//checks whether user is signin or not
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		utils.Logger.Warn("Unauthorized.User not authenticated", zap.Error(err), zap.Int64("userId", userId.(int64)))
		c.Set("response", nil)
		c.Set("message", "Unauthorized.User not authenticated")
		c.Set("error", true)
		c.Status(http.StatusUnauthorized)
		return
	}

	//Get user by ID
	user, err := dao.GetUserById(userId.(int64))
	if err != nil {
		utils.Logger.Error("Could not fetch user", zap.Error(err), zap.Int64("userId", userId.(int64)))
		c.Set("response", nil)
		c.Set("message", "Could not fetch user")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	user.Avatar = "/user/avatar/" + strconv.FormatInt(userId.(int64), 10)
	utils.Logger.Info("User Fetch successfully", zap.Error(err), zap.Int64("userId", userId.(int64)))
	c.Set("response", gin.H{"user": user})
	c.Set("message", "User Fetch successfully")
	c.Set("error", false)
	c.Status(http.StatusOK)

}

//Updates user details
func UpdateUser(c *gin.Context) {
	var req models.UpdateUserRequest
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}
	userId, exists := c.Get("userId")
	if !exists {
		utils.Logger.Warn("Unauthorized .User not authenticated", zap.Error(err), zap.Int64("userId", userId.(int64)))
		c.Set("response", nil)
		c.Set("message", "Unauthorized .User not authenticated")
		c.Set("error", true)
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Error("Invalid request body", zap.Error(err), zap.Int64("userId", userId.(int64)))
		c.Set("response", nil)
		c.Set("message", "Invalid request body")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	//Validates the user details taken from user response to update details
	err = utils.ValidateUser(req.Name, req.Mobile_No)
	if err != nil {
		utils.Logger.Error("Credentials are not validate", zap.Error(err), zap.Int64("userId", userId.(int64)))
		c.Set("response", nil)
		c.Set("message", err.Error())
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	//Updates user details
	err = dao.UpdateUserDetails(userId.(int64), req)
	if err != nil {
		utils.Logger.Error("Failed to update user", zap.Error(err), zap.Int64("userId", userId.(int64)))
		c.Set("response", nil)
		c.Set("message", "Failed to update user")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("User details updated successfully", zap.Int64("userId", userId.(int64)))
	c.Set("response", nil)
	c.Set("message", "User details updated successfully")
	c.Set("error", false)
	c.Status(http.StatusOK)

}
