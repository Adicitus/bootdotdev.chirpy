package main

import (
	"crypto/rand"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/Adicitus/bootdotdev.chirpy/frontend"
	"github.com/Adicitus/bootdotdev.chirpy/internal/database"
	"github.com/Adicitus/bootdotdev.chirpy/trie"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed sql/schema/*.sql
var migrations embed.FS

type ChirpyContext struct {
	Stats    *ApiStats
	DB       *database.Queries
	BadWords *trie.TrieNode
	TokenKey []byte
}

func reportError(w http.ResponseWriter, err error, code int) {
	data, err := json.Marshal(ChirpError{
		Error: err.Error(),
	})

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("server error"))
		return
	}

	w.WriteHeader(code)
	if os.Getenv("PLATFORM") == "dev" {
		w.Write(data)
	} else {
		w.Write([]byte(http.StatusText(code)))
	}
}

func reportResult[T any](w http.ResponseWriter, result T, code int) {
	data, err := json.Marshal(result)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("server error"))
	}

	w.WriteHeader(code)
	w.Write(data)
}

func readRequestBody[T any](r *http.Request) (v T, err error) {
	err = json.NewDecoder(r.Body).Decode(&v)
	return
}

func main() {

	godotenv.Load()

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	goose.SetBaseFS(migrations)

	dbVersion, err := goose.GetDBVersion(db)

	if err != nil {
		fmt.Printf("Error: %e\n", err)
		os.Exit(1)
	}

	fmt.Printf("Current DB version: %d\n", dbVersion)

	err = goose.Up(db, "sql/schema")

	if err != nil {
		fmt.Printf("Error: %e\n", err)
		os.Exit(1)
	}

	dbVersion, err = goose.GetDBVersion(db)

	if err != nil {
		fmt.Printf("Error: %e\n", err)
		os.Exit(1)
	}

	fmt.Printf("DB version: %d\n", dbVersion)

	cctx := new(ChirpyContext)

	if key := os.Getenv("TOKEN_KEY"); key != "" {
		cctx.TokenKey = []byte(key)
	} else {
		// Use a single-instance secret, this means all tokens will be invalidated whenever the server restarts
		b := make([]byte, 256)
		_, err := rand.Read(b)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		cctx.TokenKey = b
	}

	cctx.DB = database.New(db)

	cctx.Stats = new(ApiStats)

	cctx.BadWords = trie.NewTrie()
	cctx.BadWords.CaseInsensitive = true

	cctx.BadWords.Add("kerfuffle")
	cctx.BadWords.Add("sharbert")
	cctx.BadWords.Add("fornax")

	files := http.StripPrefix("/app", http.FileServer(http.FS(frontend.FS)))
	mux := http.NewServeMux()

	// User endpoints that can be accessed without authentication
	mux.Handle("/app/", cctx.Stats.HitsCounter(files))
	mux.HandleFunc("GET /api/chirps", handleGetChirps(cctx))
	mux.HandleFunc("GET /api/chirps/{chirpID}", handleGetChirp(cctx))
	mux.HandleFunc("POST /api/login", handleLogin(cctx))

	// User endpoints that should normally require an access token but don't when running in dev mode (boot.dev tutorial requirement)
	if os.Getenv("PLATFORM") == "dev" {
		mux.HandleFunc("POST /api/users", handleCreateUser(cctx))
		mux.HandleFunc("POST /admin/reset", handleAdminReset(cctx))
	} else {
		mux.HandleFunc("POST /api/users", secureAccess(cctx, handleCreateUser(cctx)))
		mux.HandleFunc("POST /admin/reset", secureAccess(cctx, handleAdminReset(cctx)))
	}

	// User endpoints that require refresh tokens
	mux.HandleFunc("POST /api/refresh", secureRefresh(cctx, handleTokenRefresh(cctx)))
	mux.HandleFunc("POST /api/revoke", secureRefresh(cctx, handleTokenRevoke(cctx)))

	// User endpoints requiring an Access Token
	mux.HandleFunc("GET /api/healthz", secureAccess(cctx, handleHealthz(cctx)))
	mux.HandleFunc("GET /admin/metrics", secureAccess(cctx, handleAdminMetrics(cctx)))
	mux.HandleFunc("POST /api/validate_chirp", secureAccess(cctx, handleValidateChirp(cctx)))
	mux.HandleFunc("POST /api/chirps", secureAccess(cctx, handleCreateChirp(cctx)))
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", secureAccess(cctx, handleRemoveChirp(cctx)))
	mux.HandleFunc("PUT /api/users", secureAccess(cctx, handleUpdateUser(cctx)))

	// Webhooks, security depends on the provider so should be managed in the handler
	mux.HandleFunc("POST /api/polka/webhooks", handlePolkaWebhook(cctx))

	var server http.Server

	server.Handler = mux
	server.Addr = ":8080"

	server_stop := make(chan error, 1)
	signal_stop := make(chan os.Signal, 1)

	signal.Notify(signal_stop)

	go func() {
		err := server.ListenAndServe()
		server_stop <- err
	}()

	fmt.Printf("Server listening on port 8080\n")

	for {
		select {
		case e := <-server_stop:
			fmt.Printf("Unexpected server shutdown: %s\n", e)
			os.Exit(1)
		case s := <-signal_stop:
			fmt.Printf("Signal: %s\n", s)
			switch s {
			case os.Interrupt:
				os.Exit(0)
			case os.Kill:
				os.Exit(1)
			}
		}
	}

}
