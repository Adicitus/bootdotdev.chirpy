package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Adicitus/bootdotdev.chirpy/internal/auth"
	"github.com/Adicitus/bootdotdev.chirpy/internal/database"
	"github.com/google/uuid"
)

func handleHealthz(_ *ChirpyContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("OK"))
	}
}

func handleAdminMetrics(cctx *ChirpyContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := fmt.Sprintf(`
<html>
<body>
<h1>Welcome, Chirpy Admin</h1>
<p>Chirpy has been visited %d times!</p>
</body>
</html>`,
			cctx.Stats.hits.Load())
		sr := strings.NewReader(msg)
		msg_b := make([]byte, len(msg))
		_, err := sr.Read(msg_b)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(msg_b)
	}
}

func handleAdminReset(cctx *ChirpyContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("PLATFORM") == "dev" {
			cctx.Stats.Reset()
			cctx.DB.ClearUsers(r.Context())
			w.WriteHeader(200)
		} else {
			w.WriteHeader(403)
		}
	}
}

func handleValidateChirp(cctx *ChirpyContext) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength == 0 {
			w.WriteHeader(400)
			w.Write([]byte("No chirp submitted"))
		}
		defer r.Body.Close()

		chirp, err := readRequestBody[ChirpSubmission](r)

		if err != nil {
			reportError(w, err, 400)
			return
		}

		valid, err := validateChirp(r.Context(), cctx, &chirp)

		if err != nil {
			reportError(w, err, 400)
			return
		}

		data, err := json.Marshal(valid)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Server error"))
			return
		}

		w.WriteHeader(200)
		w.Write(data)
	}
}

func handleCreateUser(cctx *ChirpyContext) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println(0)
		details, err := readRequestBody[UserDetails](r)

		if err != nil {
			reportError(w, err, 400)
			return
		}

		fmt.Println(1)
		err = auth.ValidatePassword(details.Password)

		if err != nil {
			reportError(w, err, 400)
			return
		}

		fmt.Println(2)
		user, err := cctx.DB.CreateUser(r.Context(), details.Email)

		if err != nil {
			reportError(w, err, 400)
			return
		}

		fmt.Println(3)
		fmt.Printf("Password: %s\n", details.Password)
		hash, err := auth.HashPassword(details.Password)
		fmt.Println(3.5)
		if err != nil {
			reportError(w, fmt.Errorf("Unable to process user details."), 500)
			return
		}

		fmt.Println(4)
		_, err = cctx.DB.CreateIdentity(r.Context(), database.CreateIdentityParams{
			UserID: user.ID,
			Auth:   hash,
		})

		fmt.Println(5)

		if err != nil {
			err = cctx.DB.RemoveUser(r.Context(), user.ID)
			if err != nil {
				reportError(w, fmt.Errorf("Unexpected server error."), 500)
				return
			}
			reportError(w, fmt.Errorf("Unable to process user details."), 500)
			return
		}

		reportResult(w, user, 201)
	}
}

func handleCreateChirp(cctx *ChirpyContext) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()
		chirpSubmission, err := readRequestBody[ChirpSubmission](r)

		if err != nil {
			reportError(w, err, 400)
			return
		}

		valid, err := validateChirp(r.Context(), cctx, &chirpSubmission)

		if err != nil {
			reportError(w, err, 400)
			return
		}

		chirp, err := cctx.DB.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   valid.CleanedBody,
			UserID: valid.UserID,
		})

		if err != nil {
			reportError(w, err, 400)
			return
		}

		data, err := json.Marshal(chirp)

		if err != nil {
			reportError(w, err, 500)
		}

		w.WriteHeader(201)
		w.Write(data)
	}
}

func handleGetChirps(cctx *ChirpyContext) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		chirps, err := cctx.DB.GetChirps(r.Context())

		if err != nil {
			reportError(w, err, 500)
			return
		}

		data, err := json.Marshal(chirps)

		if err != nil {
			reportError(w, err, 500)
			return
		}

		w.WriteHeader(200)
		w.Write(data)
	}
}

func handleGetChirp(cctx *ChirpyContext) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		id_s := r.PathValue("chirpID")

		if id_s == "" {
			reportError(w, fmt.Errorf("No chirp ID provided"), 404)
			return
		}

		id, err := uuid.Parse(id_s)

		if err != nil {
			reportError(w, err, 400)
			return
		}

		chirp, err := cctx.DB.GetChirp(r.Context(), id)

		if err != nil {
			reportError(w, err, 404)
			return
		}

		data, err := json.Marshal(chirp)

		if err != nil {
			reportError(w, err, 500)
			return
		}

		w.WriteHeader(200)
		w.Write(data)
	}
}

func handleLogin(cctx *ChirpyContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		details, err := readRequestBody[UserDetails](r)

		if err != nil {
			reportError(w, err, 400)
		}

		user, err := cctx.DB.GetUserByEmail(r.Context(), details.Email)

		if err != nil {
			reportError(w, err, 400)
		}

		id, err := cctx.DB.GetIdentity(r.Context(), user.ID)

		if err != nil {
			reportError(w, err, 500)
		}

		verified, err := auth.VerifyPassword(details.Password, id.Auth)

		if err != nil {
			reportError(w, err, 401)
			return
		}

		if !verified {
			reportError(w, fmt.Errorf("Invalid email/password."), 401)
			return
		}

		reportResult(w, user, 200)
	}
}
