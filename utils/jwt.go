package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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

func VerifyJwtToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("UNEXPECTED SIGNING METHOD")
		}
		return []byte(os.Getenv("JWT_SEC")), nil
	})
	if err != nil {
		return 0, errors.New("Could not parse the token.")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return 0, errors.New("Invalid Token!")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("Invalid token claims")
	}

	// email := claims["email"].(string)
	userId := int64(claims["userId"].(float64))

	return userId, nil
}
