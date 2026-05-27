package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Adicitus/bootdotdev.chirpy/internal/auth"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type LoginDetails struct {
	UserDetails
	expires_in_seconds int64
}

func secure(cctx *ChirpyContext, handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		auth_string := r.Header.Get("Authorization")

		if auth_string == "" {
			reportError(w, fmt.Errorf("Unauthorized"), 403)
			return
		}

		parts := strings.Split(auth_string, " ")

		if len(parts) != 2 {
			reportError(w, fmt.Errorf("Unauthorized"), 403)
			return
		}

		if strings.ToLower(parts[0]) != "bearer" {
			reportError(w, fmt.Errorf("Unauthorized"), 403)
			return
		}

		verified, err := auth.VerifyToken(parts[1], cctx.TokenKey)

		if err != nil {
			reportError(w, fmt.Errorf("Unauthorized"), 403)
			return
		}

		userID, err := verified.Claims.GetSubject()

		if err != nil {
			reportError(w, err, 400)
		}

		r.Header.Set("X-Chirpy-UserID", userID)

		handler(w, r)
	}
}
