package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/0x6a616e/notes/templates"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(start))
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("notes/")
	if err != nil {
		log.Println(err)
	}
	entries := []string{}
	for _, file := range files {
		entries = append(entries, file.Name())
	}
	if err = templates.Index(entries).Render(r.Context(), w); err != nil {
		log.Println(err)
	}
}

func renderFile(w http.ResponseWriter, r *http.Request) {
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /notes", index)
	mux.HandleFunc("GET /notes/{filename}", renderFile)

	return mux
}

func main() {
	server := http.Server{
		Addr:    ":8000",
		Handler: logging(newMux()),
	}
	log.Fatal(server.ListenAndServe())
}
