package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

func main() {

	var stats ApiStats
	files := http.StripPrefix("/app", http.FileServer(http.Dir("./static")))
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("OK"))
	})

	mux.Handle("/app/", stats.HitsCounter(files))

	mux.HandleFunc("GET /api/metrics", func(w http.ResponseWriter, r *http.Request) {
		msg := fmt.Sprintf("Hits: %d", stats.hits.Load())
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
	})
	mux.HandleFunc("POST /api/reset", func(w http.ResponseWriter, r *http.Request) {
		stats.Reset()
	})

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
