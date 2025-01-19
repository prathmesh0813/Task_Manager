package dao

import (
	"errors"
	"task_manager/models"
	"task_manager/utils"
	"time"

	"go.uber.org/zap"
)

// Save user in DB
func SaveUser(u *models.User) (int64, error) {
	user := User{
		Name:     u.Name,
		MobileNo: u.Mobile_No,
		Gender:   u.Gender,
		Email:    u.Email,
	}
	if err := DB.Create(&user).Error; err != nil {
		utils.Logger.Error("Failed to save user in user table", zap.Error(err), zap.Int64("userId", u.ID))
		return 0, err
	}

	u.ID = user.ID

	//Hashed Password
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		utils.Logger.Error("Failed to hashed password", zap.Error(err))
		return 0, err
	}

	login := models.Login{
		Email:    u.Email,
		Password: hashedPassword,
		UserID:   u.ID,
	}
	if err := DB.Create(&login).Error; err != nil {
		utils.Logger.Error("Failed to save user in login table", zap.Error(err), zap.Int64("userId", u.ID))
		return 0, err
	}
	utils.Logger.Info("User save successfully", zap.Int64("userId", u.ID))
	return u.ID, nil
}

// Save pair of tokens in db
func SaveToken(uid int64, user_token, refresh_token string) error {
	token := Token{
		UserToken:    user_token,
		RefreshToken: refresh_token,
		Timestamp:    time.Now(),
		UserID:       uid,
	}
	if err := DB.Create(&token).Error; err != nil {
		utils.Logger.Error("Token not save ", zap.Error(err), zap.Int64("userId", uid))
		return err
	}
	utils.Logger.Info("Token save successfully", zap.Int64("userId", uid))
	return nil
}

// validate credentials
func ValidateCredentials(u *models.Login) error {
	var login models.Login
	if err := DB.Where("email = ?", u.Email).First(&login).Error; err != nil {
		utils.Logger.Warn("Invalid credentials provided", zap.Error(err))
		return errors.New("invalid credentials")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, login.Password)
	if !passwordIsValid {
		utils.Logger.Warn("Password mismatch")
		return errors.New("invalid credentials")
	}

	utils.Logger.Info("User credentials validated successfully")
	u.ID = login.ID
	return nil
}

// Fetches user details from DB
func GetUserById(uid int64) (*models.UserResponse, error) {

	var user models.UserResponse
	result := DB.Model(User{}).Select("name, mobile_no, gender, email ").Where("id =?", uid).First(&user)
	if result.Error != nil {
		utils.Logger.Error("Failed to fetch user by id", zap.Error(result.Error), zap.Int64("userId", uid))
		return nil, result.Error
	}
	utils.Logger.Info("Fetched Task by id successfully", zap.Int64("userId", uid))
	return &user, nil
}

// delete refresh token from DB
func DeleteRefreshToken(tokenString string) error {
	var token Token
	if err := DB.Where("refresh_token = ?", tokenString).First(&token).Error; err != nil {
		utils.Logger.Error("Refresh Token not found for deletion", zap.Error(err))
		return err
	}

	if err := DB.Delete(&token).Error; err != nil {
		utils.Logger.Error("Failed to delete refresh token", zap.Error(err))
		return err
	}

	utils.Logger.Info("Refresh Token deleted successfully")
	return nil
}

// Updates user details in DB
func UpdateUserDetails(uid int64, req models.UpdateUserRequest) error {

	result := DB.Model(&User{}).Where("id = ?", uid).Updates(map[string]interface{}{
		"name":      req.Name,
		"mobile_no": req.Mobile_No,
	})
	if result.Error != nil {
		utils.Logger.Error("User not updated ", zap.Error(result.Error), zap.Int64("userid", uid))
		return result.Error
	}
	utils.Logger.Info("User updated successfully", zap.Int64("UserId", uid))
	return nil

}

// Fetches the user from DB
func GetUserByIdPassChng(uid int64) (*Login, error) {
	var login Login
	if err := DB.Where("id = ?", uid).First(&login).Error; err != nil {
		utils.Logger.Error("User not found", zap.Error(err), zap.Int64("userId", uid))
		return nil, err
	}
	utils.Logger.Info("User fetch successfully", zap.Int64("userId", uid))
	return &login, nil
}

//Update password in DB
func UpdatePassById(uid int64, password string) error {
	if err := DB.Model(&Login{}).Where("user_id = ?", uid).Update("password", password).Error; err != nil {
		utils.Logger.Error("Password not updated", zap.Error(err), zap.Int64("userid", uid))
		return err
	}
	utils.Logger.Info("Password updated successfully", zap.Int64("UserId", uid))
	return nil
}

//Delete all the tokens from DB except the token which user is login
func DeleteTokenById(uid int64, tokenString string) error {

	var token Token

	result := DB.Where("user_id = ? AND user_token != ?", uid, tokenString).Delete(&token)

	if result.RowsAffected == 0 {
		return nil
	}

	if result.Error != nil {
		utils.Logger.Error("User token not deleted ", zap.Error(result.Error), zap.Int64("userid", uid))
		return result.Error
	}
	utils.Logger.Info("User token deleted", zap.Int64("UserId", uid))
	return nil
}
