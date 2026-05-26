package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/Adicitus/bootdotdev.chirpy/internal/database"
	"github.com/Adicitus/bootdotdev.chirpy/trie"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type ChirpyContext struct {
	Stats    *ApiStats
	DB       *database.Queries
	BadWords *trie.TrieNode
}

func reportError(w http.ResponseWriter, err error, code int) {
	data, err := json.Marshal(ChirpError{
		Error: err.Error(),
	})

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("server error"))
	}

	w.WriteHeader(code)
	w.Write(data)
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
	}

	cctx := new(ChirpyContext)

	cctx.DB = database.New(db)

	cctx.Stats = new(ApiStats)

	cctx.BadWords = trie.NewTrie()
	cctx.BadWords.CaseInsensitive = true

	cctx.BadWords.Add("kerfuffle")
	cctx.BadWords.Add("sharbert")
	cctx.BadWords.Add("fornax")

	files := http.StripPrefix("/app", http.FileServer(http.Dir("./static")))
	mux := http.NewServeMux()

	mux.Handle("/app/", cctx.Stats.HitsCounter(files))
	mux.HandleFunc("GET /api/healthz", handleHealthz(cctx))
	mux.HandleFunc("GET /admin/metrics", handleAdminMetrics(cctx))
	mux.HandleFunc("POST /admin/reset", handleAdminReset(cctx))
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp(cctx))
	mux.HandleFunc("POST /api/users", handleCreateUser(cctx))
	mux.HandleFunc("POST /api/chirps", handleCreateChirp(cctx))
	mux.HandleFunc("GET /api/chirps", handleGetChirps(cctx))
	mux.HandleFunc("GET /api/chirps/{chirpID}", handleGetChirp(cctx))
	mux.HandleFunc("POST /api/login", handleLogin(cctx))

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
