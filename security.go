/*
security.go

Contains functionality related to logins and security that are not handlers and are not covered by the auth or database modules.
*/

package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Adicitus/bootdotdev.chirpy/internal/auth"
	"github.com/Adicitus/bootdotdev.chirpy/internal/database"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type LoginDetails struct {
	UserDetails
}

type UserAuthDetails struct {
	database.User
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func getAuthorizationToken(r *http.Request) (token string, err error) {
	authorizationString := r.Header.Get("Authorization")

	if authorizationString == "" {
		return "", fmt.Errorf("No Authorization header")
	}

	parts := strings.Split(authorizationString, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("Malformed Authorization header")
	}

	return parts[1], nil
}

func secure(cctx *ChirpyContext, handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		token, err := getAuthorizationToken(r)

		if err != nil {
			reportError(w, err, 401)
			return
		}

		verified, err := auth.VerifyToken(token, cctx.TokenKey)

		if err != nil {
			reportError(w, err, 401)
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
