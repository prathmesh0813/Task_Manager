package utils

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

// ValidateDetails validates name, email, mobile, gender, and password.
func ValidateDetails(name, email, mobile, gender, password string) error {
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
		return errors.New("mobile number must be exactly 10 digits and contain only numbers")
	}
	// Validate gender
	gender = strings.ToLower(gender)
	if gender != "male" && gender != "female" && gender != "other" {
		return errors.New("gender must be male, female or other")
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
func ValidateUser(name string, mobileno string) error {
	// Validate name
	if len(name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}
	if !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(name) {
		return errors.New("name must contain only alphabets and spaces")
	}

	// Validate mobile

	if len(mobileno) != 10 || !regexp.MustCompile(`^\d{10}$`).Match([]byte(mobileno)) {
		return errors.New("mobile number must be exactly 10 digits and contain only numbers")
	}

	return nil
}

// ValidatePassword checks if both old and new passwords meet complexity rules
func ValidatePassword(oldPassword, newPassword string) error {
	// Validate old password
	if err := checkPasswordComplexity(oldPassword); err != nil {
		return errors.New("old password: " + err.Error())
	}

	// Validate new password
	if err := checkPasswordComplexity(newPassword); err != nil {
		return errors.New("new password: " + err.Error())
	}

	return nil // Both passwords are valid
}

// checkPasswordComplexity enforces password strength rules
func checkPasswordComplexity(password string) error {
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

	// Ensure password is at least 8 characters long
	if len(password) < 8 {
		return errors.New("must be at least 8 characters long")
	}

	// Ensure all complexity rules are met
	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return errors.New("must contain at least one uppercase letter, one lowercase letter, one digit, and one special character")
	}

	return nil // Password is valid
}

func ValidateLoginDetails(email, password string) error {
	// Validate email
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
		return errors.New("invalid email format")
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
