package dao

import (
	"errors"
	"task_manager/models"
	"task_manager/utils"
	"time"

	"go.uber.org/zap"
)

// Save user in DB
func (u *models.User) Save() (int64, error) {
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
func (u *models.User) SaveToken(uid int64, user_token, refresh_token string) error {
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
func (u *models.Login) ValidateCredentials() error {
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
