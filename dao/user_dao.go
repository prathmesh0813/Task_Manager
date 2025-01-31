package dao

import (
	"errors"
	"task_manager/models"
	"task_manager/utils"
	"time"

	"gorm.io/gorm"
)

// Save user in DB
func SaveUser(db *gorm.DB, u *models.User) (int64, string, string, error) {
	// Note the use of tx as the database handle once you are within a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return 0, "", "", err
	}

	//hashing password
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		tx.Rollback()
		return 0, "", "", err
	}

	//saving user in db
	user := User{
		Name:     u.Name,
		MobileNo: u.Mobile_No,
		Gender:   u.Gender,
		Email:    u.Email,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return 0, "", "", err
	}

	u.ID = user.ID

	//Generate user token
	userToken, refreshToken, err := utils.GenerateTokens(u.ID)
	if err != nil {
		tx.Rollback()
		return 0, "", "", err
	}

	//saving email and hashed password in db
	login := models.Login{
		Email:    u.Email,
		Password: hashedPassword,
		UserID:   u.ID,
	}

	if err := tx.Create(&login).Error; err != nil {
		tx.Rollback()
		return 0, "", "", err
	}

	//saving token in db
	token := Token{
		UserToken:    userToken,
		RefreshToken: refreshToken,
		Timestamp:    time.Now(),
		UserID:       u.ID,
	}

	if err := tx.Create(&token).Error; err != nil {
		tx.Rollback()
		return 0, "", "", err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, "", "", err
	}

	return u.ID, userToken, refreshToken, nil
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
		return err
	}

	return nil
}

// validate credentials
func ValidateCredentials(u *models.Login) error {
	var login models.Login
	if err := DB.Where("email = ?", u.Email).First(&login).Error; err != nil {
		return errors.New("invalid credentials")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, login.Password)
	if !passwordIsValid {
		return errors.New("invalid credentials")
	}

	u.ID = login.ID
	return nil
}

// Fetches user details from DB
func GetUserById(uid int64) (*models.UserResponse, error) {

	var user models.UserResponse
	result := DB.Model(User{}).Select("name, mobile_no, gender, email ").Where("id =?", uid).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// delete refresh token from DB
func DeleteRefreshToken(tokenString string) error {
	var token Token
	if err := DB.Where("refresh_token = ?", tokenString).First(&token).Error; err != nil {
		return err
	}

	if err := DB.Delete(&token).Error; err != nil {
		return err
	}

	return nil
}

// Updates user details in DB
func UpdateUserDetails(uid int64, req models.UpdateUserRequest) error {

	result := DB.Model(&User{}).Where("id = ?", uid).Updates(map[string]interface{}{
		"name":      req.Name,
		"mobile_no": req.Mobile_No,
	})
	if result.Error != nil {
		return result.Error
	}

	return nil

}

// Fetches the user from DB
func GetUserByIdPassChng(uid int64) (*Login, error) {
	var login Login
	if err := DB.Where("id = ?", uid).First(&login).Error; err != nil {
		return nil, err
	}

	return &login, nil
}

// Update password in DB
func UpdatePassById(uid int64, password string) error {
	if err := DB.Model(&Login{}).Where("user_id = ?", uid).Update("password", password).Error; err != nil {
		return err
	}

	return nil
}

// Delete all the tokens from DB except the token which user is login
func DeleteTokenById(uid int64, tokenString string) error {
	var token Token

	result := DB.Where("user_id = ? AND user_token != ?", uid, tokenString).Delete(&token)

	if result.RowsAffected == 0 {
		return nil
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// signout all the users
func SignOutAllUsers(uid int64) error {
	var token Token

	result := DB.Where("user_id = ?", uid).Delete(&token)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Signout single user
func DeleteToken(tokenString string) error {
	var token Token

	if err := DB.Where("user_token = ? ", tokenString).First(&token).Error; err != nil {
		return err
	}

	if err := DB.Delete(&token).Error; err != nil {
		return err
	}

	return nil
}

// Fetches the user from Token DB
func GetUserByIdFromTokenTable(uid int64) error {
	var token Token
	if err := DB.Where("id = ?", uid).First(&token).Error; err != nil {
		return err
	}

	return nil
}
