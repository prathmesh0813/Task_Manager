package controller

import (
	"net/http"

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
		c.Set("message", "cannot parsed the requested data")
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
	uid, err := user.Save()
	if err != nil {
		utils.Logger.Error("Could not save user. EmailId already existed", zap.Error(err), zap.Int64("userId", user.ID))
		c.Set("response", nil)
		c.Set("message", "Could not save user. EmailId already existed")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return

	}

	//Generate user token
	userToken, err := utils.GenerateJwtToken(uid)
	if err != nil {
		utils.Logger.Error("could not generate the user token", zap.Error(err), zap.Int64("userId", user.ID))
		c.Set("response", nil)
		c.Set("message", "could not generate the user token")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	//Generate Refresh token
	refreshToken, err := utils.GenerateRefreshToken(uid)
	if err != nil {
		utils.Logger.Error("could not generate the refresh token", zap.Error(err), zap.Int64("userId", user.ID))
		c.Set("response", nil)
		c.Set("message", "could not generate the refresh token")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	//Save tokens
	err = user.SaveToken(uid, userToken, refreshToken)
	if err != nil {
		utils.Logger.Error("could not save the token", zap.Error(err), zap.Int64("userId", user.ID))
		c.Set("response", nil)
		c.Set("message", "could not save the token")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("User save successfully", zap.Int64("userId", uid))
	c.Set("response", gin.H{"refresh_token": refreshToken, "user_token": userToken})
	c.Set("message", "User save successfully")
	c.Set("error", false)
	c.Status(http.StatusCreated)
}

// user sign in
func SignIn(c *gin.Context) {
	var login models.Login
	var user models.User

	err := c.ShouldBindJSON(&login)
	if err != nil {
		utils.Logger.Warn("Failed to parse login request", zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "could not parse the request data")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	//validate credentials
	err = login.ValidateCredentials()
	if err != nil {
		utils.Logger.Warn("Authentication failed", zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "could not authenticate user")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	//generate user token
	userToken, err := utils.GenerateJwtToken(login.ID)
	if err != nil {
		utils.Logger.Error("Failed to generate user token", zap.Int64("userId", login.ID), zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "could not generate the token")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	//generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(login.ID)
	if err != nil {
		utils.Logger.Error("Failed to generate refresh token", zap.Int64("userId", login.ID), zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "could not generate the token")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	//save token in db
	err = user.SaveToken(login.ID, userToken, refreshToken)
	if err != nil {
		utils.Logger.Error("Failed to save token", zap.Int64("userId", login.ID), zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "could not save the token")
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
