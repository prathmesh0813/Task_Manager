package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// ValidateDetails validates name, email, mobile, gender, and password.
func ValidateDetails(name, email string, mobile, gender, password string) error {
	// Validate name
	if len(name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}
	if !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(name) {
		return errors.New("name must contain only alphabets and spaces")
	}
	// Validate email
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
		return errors.New("invalid email format")
	}
	// Validate mobile
	if len(mobile) != 10 || !regexp.MustCompile(`^\d{10}$`).Match([]byte(mobile)) {
		return errors.New("mobile number must be 10 digits")
	}
	// Validate gender
	gender = strings.ToLower(gender)
	if gender != "male" && gender != "female" && gender != "other" {
		return errors.New("gender must be 'male', 'female', or 'other'")
	}
	// Validate password
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {

		return errors.New("password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character")
	}
	return nil

}

// Validated user details
func ValidateUser(name string, mobileno int64) error {
	// Validate name
	if len(name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}
	if !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(name) {
		return errors.New("name must contain only alphabets and spaces")
	}

	// Validate mobile

	mobileStr := strconv.FormatInt(mobileno, 10)
	if len(mobileStr) != 10 || !regexp.MustCompile(`^\d{10}$`).Match([]byte(mobileStr)) {
		return errors.New("mobile number must be 10 digits")
	}

	return nil
}

// Validate password
func ValidatePassword(password string) error {
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {

		return errors.New("password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character")
	}
	return nil
}
