package utils

import "golang.org/x/crypto/bcrypt"

// Hash the plain string
func CreateHash(plainString string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainString), 7)
	return string(bytes), err
}

// checks the plain text password with hash password
func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
