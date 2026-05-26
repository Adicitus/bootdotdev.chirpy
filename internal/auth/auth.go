package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func ValidatePassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("Password too short (must be 8 characters or longer)")
	}

	return nil
}

func HashPassword(password string) (hash string, err error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func VerifyPassword(password, hash string) (correct bool, err error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
