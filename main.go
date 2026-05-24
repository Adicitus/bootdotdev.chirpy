package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/Adicitus/bootdotdev.chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type ChirpyContext struct {
	Stats *ApiStats
	DB    *database.Queries
}

func reportError(w http.ResponseWriter, err error) {
	data, err := json.Marshal(ChirpError{
		Error: err.Error(),
	})

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("server error"))
	}

	w.WriteHeader(400)
	w.Write(data)
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
	files := http.StripPrefix("/app", http.FileServer(http.Dir("./static")))
	mux := http.NewServeMux()

	mux.Handle("/app/", cctx.Stats.HitsCounter(files))
	mux.HandleFunc("GET /api/healthz", handleHealthz(cctx))
	mux.HandleFunc("GET /admin/metrics", handleAdminMetrics(cctx))
	mux.HandleFunc("POST /admin/reset", handleAdminReset(cctx))
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp(cctx))
	mux.HandleFunc("POST /api/users", handleCreateUser(cctx))

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

	select {
	case e := <-server_stop:
		fmt.Printf("Unexpected server shutdown: %s\n", e)
	case s := <-signal_stop:
		fmt.Printf("Signal: %s\n", s)
	}

}
