package auth_test

import (
	"testing"

	"github.com/Adicitus/bootdotdev.chirpy/internal/auth"
)

func TestHashPassword(t *testing.T) {
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
