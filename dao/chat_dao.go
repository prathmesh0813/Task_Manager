package dao

import (
	"context"
	"errors"
	"task_manager/models"

	"gorm.io/gorm"
)

//checks whether given email id's present in (SQL)db or not
func CheckEmailPresent(emails []string) (error, []string) {
	var errFound error

	missingUsers := []string{}

	for _, email := range emails {
		var user User
		err := DB.Where("email = ?", email).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			missingUsers = append(missingUsers, email)
			errFound = gorm.ErrRecordNotFound
		} else if err != nil {
			return err, nil
		}
	}

	if errFound != nil {
		return errFound, missingUsers
	}

	return nil, nil
}

// InserChatHead inserts chat head into MongoDB
func InserChatHead(email models.Emails) error {
	_, err := ChatCollection.InsertOne(context.TODO(), email)
	if err != nil {
		return err
	}
	return nil
}
