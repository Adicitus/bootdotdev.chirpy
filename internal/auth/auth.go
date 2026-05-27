package auth

import (
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ValidatePassword(password string) error {
	if len(password) < 1 {
		return fmt.Errorf("Password too short (must be 6 characters or longer)")
	}

	return nil
}

func HashPassword(password string) (hash string, err error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func VerifyPassword(password, hash string) (correct bool, err error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func CreateToken(userID uuid.UUID, valididty time.Duration, secret []byte) (string, error) {
	now := time.Now()
	exp := now.Add(valididty)

	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": exp.Unix(),
			"iat": now.Unix(),
			"nbf": now.Unix(),
			"sub": userID.String(),
			"iss": "chirpy-access",
		})

	return t.SignedString(secret)
}

func VerifyToken(tokenString string, secret []byte) (verified *jwt.Token, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, err
}
