package dao

import (
	"errors"

	"gorm.io/gorm"
)

func CheckEmailPresent(emails []string) (error,[]string){
	var errFound error

	missingUsers := []string{}

	for _,email := range emails{
		var user User
		err := DB.Where("email = ?",email).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound){
			missingUsers = append(missingUsers, email)
			errFound = gorm.ErrRecordNotFound
		}else if err!=nil{
			return err, nil
		}
	}

	if errFound!=nil{
		return errFound,missingUsers
	}

	return nil,nil
}