/*
security.go

Contains functionality related to logins and security that are not handlers and are not covered by the auth or database modules.
*/

package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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

type AuthorizationHeader struct {
	method string
	token  string
}

func getAuthorizationHeader(r *http.Request) (authHeader AuthorizationHeader, err error) {
	authorizationString := r.Header.Get("Authorization")

	if authorizationString == "" {
		return AuthorizationHeader{}, fmt.Errorf("No Authorization header")
	}

	parts := strings.Split(authorizationString, " ")

	if len(parts) != 2 {
		return AuthorizationHeader{}, fmt.Errorf("Malformed Authorization header")
	}

	authHeader = AuthorizationHeader{
		method: strings.ToLower(parts[0]),
		token:  parts[1],
	}

	return
}

/*
Wraps the given handler function in a middleware that ensures the "Authorization"
header of the request contains a valid access token before calling the handler.

# If the Authorization header is missing or cannot be valdated, the handler will not be called.

If the Authorization header is contains a verified token, the user ID of the associated
user will be added to the request headers as "X-Chirpy-UserID".
*/
func secureAccess(cctx *ChirpyContext, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header, err := getAuthorizationHeader(r)

		if err != nil {
			reportError(w, err, 401)
			return
		}

		if header.method != "bearer" {
			reportError(w, fmt.Errorf("Invalid authorization method"), 401)
			return
		}

		verified, err := auth.VerifyToken(header.token, cctx.TokenKey)

		if err != nil {
			reportError(w, err, 401)
			return
		}

		userID, err := verified.Claims.GetSubject()

		if err != nil {
			reportError(w, err, 400)
			return
		}

		r.Header.Set("X-Chirpy-UserID", userID)

		handler(w, r)
	}
}

/*
Wraps the given handler function in a middleware that ensures the "Authorization"
header of the request contains a valid refresh token before calling the handler.

If the Authorization header is missing or cannot be validated, the handler will not be called.

If the Authorization header is contains a verified token, the user ID of the associated
user will be added to the request headers as "X-Chirpy-UserID".
*/
func secureRefresh(cctx *ChirpyContext, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header, err := getAuthorizationHeader(r)

		if err != nil {
			reportError(w, err, 401)
			return
		}

		if header.method != "bearer" {
			reportError(w, fmt.Errorf("Invalid authorization method"), 401)
			return
		}

		token, err := cctx.DB.GetToken(r.Context(), header.token)

		if err != nil {
			reportError(w, err, 401)
			return
		}

		if token.RevokedAt.Valid {
			reportError(w, fmt.Errorf("Token revoked"), 401)
			return
		}

		if token.ExpiresAt.Before(time.Now()) {
			reportError(w, fmt.Errorf("Token expired"), 401)
			return
		}

		r.Header.Set("X-Chirpy-UserID", token.UserID.String())

		handler(w, r)
	}
}
