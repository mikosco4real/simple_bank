package util

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateUuid() string {
	return uuid.New().String()
}

func HashPassword(s string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost) // bycrypt.DefaultCost is 10. Used to set complexity
    if err != nil {
        return "", err
    }
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, newPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(newPassword))
    return err == nil
} 