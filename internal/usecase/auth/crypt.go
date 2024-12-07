package auth

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

func (h BcryptHasher) GenerateFromPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (h BcryptHasher) CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type GoogleUUID struct{}

func (u GoogleUUID) New() string {
	return uuid.New().String()
}
