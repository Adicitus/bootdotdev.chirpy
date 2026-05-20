package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

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

	stats := new(ApiStats)
	files := http.StripPrefix("/app", http.FileServer(http.Dir("./static")))
	mux := http.NewServeMux()

	mux.Handle("/app/", stats.HitsCounter(files))
	mux.HandleFunc("GET /api/healthz", handleHealthz(stats))
	mux.HandleFunc("GET /admin/metrics", handleAdminMetrics(stats))
	mux.HandleFunc("POST /admin/reset", handleAdminReset(stats))
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp(stats))

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
