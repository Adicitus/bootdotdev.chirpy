package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ChirpSubmission struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type ChirpError struct {
	Error string `json:"error"`
}

type ChirpValid struct {
	CleanedBody string    `json:"cleaned_body"`
	UserID      uuid.UUID `json:"user_id"`
}

func validateChirp(ctx context.Context, cctx *ChirpyContext, chirp *ChirpSubmission) (v ChirpValid, err error) {

	// Chirps must have a valid user ID
	if _, err = cctx.DB.GetUser(ctx, chirp.UserID); err != nil {
		return
	}

	v.UserID = chirp.UserID

	l := len(chirp.Body)

	// Chirps cannot be empty
	if l == 0 {
		err = fmt.Errorf("Empty chirp!")
		return
	}

	// Chirps cannot be more than 140 characters
	if l > 140 {
		err = fmt.Errorf("Chirp is too long")
		return
	}

	v.CleanedBody, err = cctx.BadWords.Replace(chirp.Body, "****")

	return
}
