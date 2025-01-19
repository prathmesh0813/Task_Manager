package dao

import (
	"errors"

	"gorm.io/gorm"
)

//Fetch the avatar 
func ReadAvatar(uid int64) (*Avatar, error) {

	var avatar Avatar
	if err := DB.Where("user_id = ?", uid).First(&avatar).Error; err != nil {
		return nil, err
	}
	return &avatar, nil
}

//Update the avatar DB
func UpdateAvatar(uid int64, content []byte) error {
	var existingAvatar Avatar
	result := DB.Where("user_id = ?", uid).First(&existingAvatar)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	if result.RowsAffected > 0 {
		existingAvatar.Data = content
		if err := DB.Save(&existingAvatar).Error; err != nil {
			return err
		}
	} else {
		newAvatar := Avatar{
			Data:   content,
			UserID: uid,
		}
		if err := DB.Create(&newAvatar).Error; err != nil {
			return err
		}
	}
	return nil
}

//save avatar in DB
func SaveAvatar(uid int64, content []byte, fileName string) error {

	avatar := Avatar{
		Data:   content,
		UserID: uid,
		Name:   fileName,
	}
	err := DB.Create(&avatar)
	if err != nil {
		return err.Error
	}

	return nil
}

//delete user avatar
func DeleteAvatar(uid int64) error {
	var avatar Avatar

	result := DB.Where("user_id = ?", uid).Delete(&avatar)
	if result.RowsAffected == 0 {
		return nil
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}
