package controller

import (
	"net/http"
	"strconv"
	"strings"

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

// Updates user details
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

// Update password
func UpdatePassword(c *gin.Context) {
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
		c.Set("response", nil)
		c.Set("message", "Unauthorized.User not authenticated")
		c.Set("error", true)
		c.Status(http.StatusUnauthorized)
		return
	}
	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Error("Invalid request body", zap.Error(err), zap.Int64("userId", userIdFromToken.(int64)))
		c.Set("response", nil)
		c.Set("message", "Invalid request body")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	//Validate whether enter password is in correct format or not
	err = utils.ValidatePassword(req.NewPassword)
	if err != nil {
		utils.Logger.Error("Credentials are not validate", zap.Error(err), zap.Int64("userId", userIdFromToken.(int64)))
		c.Set("response", nil)
		c.Set("message", err.Error())
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	//Check whether user is present or not in db to chng password
	user, err := dao.GetUserByIdPassChng(userIdFromToken.(int64))
	if err != nil {
		utils.Logger.Error("User not found", zap.Error(err), zap.Int64("userId", userIdFromToken.(int64)))
		c.Set("response", nil)
		c.Set("message", "User not found")
		c.Set("error", true)
		c.Status(http.StatusNotFound)
		return
	}

	//checks the user enter old password is correct with password hash which is present in DB
	passwordIsValid := utils.CheckPasswordHash(req.OldPassword, user.Password)
	if !passwordIsValid {
		utils.Logger.Error("Incorrect old password", zap.Error(err), zap.Int64("userId", userIdFromToken.(int64)))
		c.Set("response", nil)
		c.Set("message", "Incorrect old password")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	//Hash the new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.Logger.Error("Failed to hashed password", zap.Error(err), zap.Int64("userId", userIdFromToken.(int64)))
		c.Set("response", nil)
		c.Set("message", "Failed to hashed password")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	//Update password by userId
	err = dao.UpdatePassById(userIdFromToken.(int64), hashedPassword)
	if err != nil {
		utils.Logger.Error("Failed to update password", zap.Error(err), zap.Int64("userId", userIdFromToken.(int64)))
		c.Set("response", nil)
		c.Set("message", "Failed to update password")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	//Signout all other user except the user which updates the password
	err = dao.DeleteTokenById(userIdFromToken.(int64), token)
	if err != nil {
		utils.Logger.Error("Failed to delete token", zap.Error(err), zap.Int64("userId", userIdFromToken.(int64)))
		c.Set("response", nil)
		c.Set("message", "Failed to delete token")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("Password Updated Successfully.", zap.Int64("userId", userIdFromToken.(int64)))
	c.Set("response", nil)
	c.Set("message", "Password Updated Successfully.")
	c.Set("error", false)
	c.Status(http.StatusOK)

}

// Signout user
func SignOut(c *gin.Context) {
	tokenString := strings.TrimSpace(c.GetHeader("Authorization"))
	if tokenString == "" {
		utils.Logger.Error("Token not provided")
		c.Set("response", nil)
		c.Set("message", "Token not provided")
		c.Set("error", true)
		c.Status(http.StatusUnauthorized)
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

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	//takes query parameters
	allParam := c.DefaultQuery("all", "false")
	all, err := strconv.ParseBool(allParam)
	//Check whether user is present or not in db to chng password

	if err != nil {
		utils.Logger.Error("Invalid query parameter for 'all'. It must be true or false", zap.Error(err), zap.Int64("userId", userId.(int64)))
		c.Set("response", nil)
		c.Set("message", "Invalid query parameter for 'all'. It must be true or false")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return

	}

	//checks whether user is allready signout or not
	err = dao.GetUserByIdFromTokenTable(userId.(int64))
	if err != nil {
		utils.Logger.Error("User not found allready logout", zap.Error(err), zap.Int64("userId", userId.(int64)))
		c.Set("response", nil)
		c.Set("message", "User not found allready logout")
		c.Set("error", true)
		c.Status(http.StatusNotFound)
		return
	}
	//if all query param value is true

	if all {
		//signout all users
		err = dao.SignOutAllUsers(userId.(int64))
		if err != nil {
			utils.Logger.Error("Failed to signout from all devices", zap.Error(err), zap.Int64("userId", userId.(int64)))
			c.Set("response", nil)
			c.Set("message", "Failed to signout from all devices")
			c.Set("error", true)
			c.Status(http.StatusBadRequest)
		}
		utils.Logger.Info("Signout from all devices is successfully", zap.Int64("userId", userId.(int64)))
		c.Set("response", nil)
		c.Set("message", "Signout from all devices is successfully")
		c.Set("error", false)
		c.Status(http.StatusOK)
	} else {
		//Sign out single user if all query param value is false
		err = dao.DeleteToken(tokenString)
		if err != nil {
			utils.Logger.Error("Failed to signout", zap.Error(err), zap.Int64("userId", userId.(int64)))
			c.Set("response", nil)
			c.Set("message", "Failed to signout")
			c.Set("error", true)
			c.Status(http.StatusBadRequest)
			return
		}

		utils.Logger.Info("User sign out successfully", zap.Int64("userId", userId.(int64)))
		c.Set("response", nil)
		c.Set("message", "User sign out successfully")
		c.Set("error", false)
		c.Status(http.StatusOK)
	}
}
