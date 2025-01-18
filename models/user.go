package models

import (
	dao "task_manager/DAO"
	"task_manager/utils"
	"time"

	"go.uber.org/zap"
)


//User Struct for request body
type User struct {
	ID               int64
	Name             string `json:"name"`
	Mobile_No        int64  `json:"mob_no"`
	Gender           string `json:"gender"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Confirm_Password string `json:"confirm_password"`
}

//Save user in DB
func (u *User) Save() (int64, error) {
	newUser := dao.User{
		Name:     u.Name,
		MobileNo: u.Mobile_No,
		Gender:   u.Gender,
		Email:    u.Email,
	}
	if err := dao.DB.Create(&newUser).Error; err != nil {
		utils.Logger.Error("Failed to save user in user table", zap.Error(err), zap.Int64("userId", u.ID))
		return 0, err
	}


	//Hashed Password
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		utils.Logger.Error("Failed to hashed password", zap.Error(err))
		return 0, err
	}

	login := dao.Login{
		Email:    u.Email,
		Password: hashedPassword,
		UserID: newUser.ID,
	}
	if err := dao.DB.Create(&login).Error; err != nil {
		utils.Logger.Error("Failed to save user in login table", zap.Error(err), zap.Int64("userId", newUser.ID))
		return 0, err
	}
	utils.Logger.Info("User save successfully", zap.Int64("userId", newUser.ID))
	return newUser.ID, nil
}


//Save pair of tokens in db
func (u *User) SaveToken(uid int64, user_token, refresh_token string) error {
	token := dao.Token{
		UserToken:    user_token,
		RefreshToken: refresh_token,
		Timestamp:    time.Now(),
		UserID:       uid,
	}
	if err := dao.DB.Create(&token).Error; err != nil {
		utils.Logger.Error("Token not save ", zap.Error(err), zap.Int64("userId", uid))
		return err
	}
	utils.Logger.Info("Token save successfully", zap.Int64("userId", uid))
	return nil
}
