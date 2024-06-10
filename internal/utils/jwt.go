package utils

import (
	"assignment/config"
	"assignment/internal/errors"
	"assignment/internal/models"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(user *models.User) (string, *errors.Error) {
	cfg := config.GetConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.AppKey))
	if err != nil {
		return "", errors.New(errors.ErrTokenGeneration)
	}
	return tokenString, nil
}

func DecodeToken(tokenString string) (map[string]interface{}, *errors.Error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig().AppKey), nil
	})

	if err != nil {
		return nil, errors.New(errors.ErrTokenGeneration)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New(errors.ErrTokenGeneration)
}
