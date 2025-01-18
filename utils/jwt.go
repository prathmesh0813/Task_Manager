package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// const secretKey = "superscrete"
// const secretKey2 = "supertopscrete"

// generates user token
func GenerateJwtToken(userId int64) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Minute * 5).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SEC")))
}

// generate refresh token
func GenerateRefreshToken(userId int64) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_REF_SEC")))
}
