package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// Generate pair of tokens
func GenerateTokens(userId int64) (string, string, error) {
	// Default expiration durations
	const defaultUserExpDuration = "2h"
	const defaultRefreshExpDuration = "4h"

	// Fetch user token exp duration from the environment or use default
	userExpDuration := os.Getenv("JWT_EXP_DURATION")
	if userExpDuration == "" {
		userExpDuration = defaultUserExpDuration
	}

	userDuration, _ := time.ParseDuration(userExpDuration)

	// Fetch refresh token exp duration from the environment or use default
	refreshExpDuration := os.Getenv("REF_EXP_DURATION")
	if refreshExpDuration == "" {
		refreshExpDuration = defaultRefreshExpDuration
	}

	refreshDuration, _ := time.ParseDuration(refreshExpDuration)

	// Calculate exp times for user token and refresh
	userExpTime := time.Now().Add(userDuration).Unix()
	refreshExpTime := time.Now().Add(refreshDuration).Unix()

	// Fetch secrets from environment
	userSecret := os.Getenv("JWT_SEC")
	if userSecret == "" {
		return "", "", errors.New("JWT_SEC is not set")
	}
	// Fetch secrets from environment
	refreshSecret := os.Getenv("JWT_REF_SEC")
	if refreshSecret == "" {
		return "", "", errors.New("JWT_REF_SEC is not set")
	}

	// Generate user token
	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    userExpTime,
	})
	signedUserToken, _ := userToken.SignedString([]byte(userSecret))

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    refreshExpTime,
	})
	signedRefreshToken, _ := refreshToken.SignedString([]byte(refreshSecret))

	return signedUserToken, signedRefreshToken, nil
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

// verify refresh token
func VerifyRefreshToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			Logger.Error("Unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}

		return []byte(os.Getenv("JWT_REF_SEC")), nil
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

	Logger.Info("Refresh Token verified successfully", zap.Int64("userId", int64(userId)))
	return int64(userId), nil
}
