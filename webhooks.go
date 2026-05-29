/*
webhooks.go

Contains functions to generate handler functions for webhook http endpoints.
*/

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
)

// Basic format of a Webhook request body (JSON):
type WebhookRequest[T any] struct {
	Event string `json:"event"`
	Data  T      `json:"data"`
}

type WebhookPolkaData struct {
	UserID uuid.UUID `json:"user_id"`
}

func handlePolkaWebhook(cctx *ChirpyContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		polkaKey := os.Getenv("POLKA_KEY")

		if polkaKey == "" {
			reportError(w, fmt.Errorf("No Polka key set"), 500)
			return
		}

		header, err := getAuthorizationHeader(r)

		if err != nil {
			reportError(w, err, 401)
			return
		}

		if header.method != "apikey" {
			reportError(w, fmt.Errorf("Invalid method"), 401)
			return
		}

		if header.token != polkaKey {
			reportError(w, fmt.Errorf("Invalid Polka key"), 401)
			return
		}

		details, err := readRequestBody[WebhookRequest[WebhookPolkaData]](r)

		if err != nil {
			reportError(w, err, 400)
			return
		}

		if details.Event != "user.upgraded" {
			w.WriteHeader(204)
			return
		}

		_, err = cctx.DB.GetUser(r.Context(), details.Data.UserID)

		if err != nil {
			reportError(w, err, 404)
		}

		_, err = cctx.DB.ActivateMembership(r.Context(), details.Data.UserID)

		if err != nil {
			reportError(w, err, 500)
			return
		}

		w.WriteHeader(204)
	}
}
