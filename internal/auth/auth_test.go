package auth_test

import (
	"testing"
	"time"

	"github.com/Adicitus/bootdotdev.chirpy/internal/auth"
	"github.com/google/uuid"
)

func TestPasswords(t *testing.T) {
	password := "123456"
	t.Logf("Creating a hash of '%s'...", password)
	hash, err := auth.HashPassword(password)

	if err != nil {
		t.Logf("Error: %s", err)
		t.Fail()
	}

	t.Logf("Created hash: %s", hash)

	t.Logf("Verifying password '%s' against hash '%s'...", password, hash)

	verified, err := auth.VerifyPassword(password, hash)

	if err != nil {
		t.Logf("Error: %s", err)
		t.Fail()
	}

	if !verified {
		t.Log("Failed to verify!")
		t.Fail()
	}

	verified, err = auth.VerifyPassword("incorrect password", hash)

	if err != nil {
		t.Logf("Error: %s", err)
		t.Fail()
	}

	if verified {
		t.Log("Verified incorrect password!")
		t.Fail()
	}
}

func TestTokens(t *testing.T) {
	secret1 := []byte("Ssh! This is a secret!")
	secret2 := []byte("Not a secret!")
	id := uuid.New()

	tokenString, err := auth.CreateToken(id, 5*time.Minute, secret1)

	if err != nil {
		t.Logf("failed to create token string: %s", err)
		t.Fail()
	}

	verified, err := auth.VerifyToken(tokenString, secret1)

	if err != nil {
		t.Logf("Failed to verify token: %s", err)
		t.Fail()
	}

	if !verified {
		t.Log("Valid token not verified")
		t.Fail()
	}

	verified, err = auth.VerifyToken(tokenString, secret2)

	if err == nil {
		t.Logf("Expected error when verifying using the wrong key")
		t.Fail()
	}

	if verified {
		t.Log("Verified token with the wrong secret")
		t.Fail()
	}

	tokenString2, err := auth.CreateToken(id, 5*time.Millisecond, secret1)

	time.Sleep(5 * time.Millisecond)

	verified, err = auth.VerifyToken(tokenString2, secret1)

	if err == nil {
		t.Logf("Expected an error due to expired token")
		t.Fail()
	}

	if verified {
		t.Log("Verified an expired token")
		t.Fail()
	}
}
