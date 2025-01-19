package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
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

// verify user token
func VerifyJwtToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			Logger.Error("Unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}

		return []byte(os.Getenv("JWT_SEC")), nil
	})

	if err != nil {
		Logger.Error("Failed to parse token", zap.Error(err))
		return 0, errors.New("could not parse the token")
	}

	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		Logger.Error("Invalid token")
		return 0, errors.New("invalid Token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		Logger.Error("Invalid token claims")
		return 0, errors.New("invalid token claims")
	}

	userId, ok := claims["userId"].(float64)
	if !ok {
		Logger.Error("Invalid user id in token claims")
		return 0, errors.New("invalid token claims")
	}

	Logger.Info("User Token verified successfully", zap.Int64("userId", int64(userId)))
	return int64(userId), nil
}
