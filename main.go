package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", greet)

	return mux
}

func main() {
	server := http.Server{
		Addr:    ":8000",
		Handler: newMux(),
	}
	log.Fatal(server.ListenAndServe())
}
