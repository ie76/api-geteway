package utils

import (
	"assignment/internal/errors"

	"golang.org/x/crypto/bcrypt"
)

func GenerateFromPassword(password string) (string, *errors.Error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New(errors.ErrHash)
	}
	return string(hashedPassword), nil
}
