package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
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
