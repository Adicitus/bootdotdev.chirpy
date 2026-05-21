package main

import "fmt"

type Chirp struct {
	Body string `json:"body"`
}

type ChirpError struct {
	Error string `json:"error"`
}

type ChirpValid struct {
	CleanedBody string `json:"cleaned_body"`
}

func validateChirp(chirp *Chirp) error {

	l := len(chirp.Body)

	if l == 0 {
		return fmt.Errorf("Empty chirp!")
	}

	if l > 140 {
		return fmt.Errorf("Chirp is too long")
	}

	return nil
}
