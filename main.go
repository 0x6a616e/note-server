package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"

	"github.com/0x6a616e/notes/internal"
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
	entries := []internal.File{}
	for _, file := range files {
		entries = append(entries, internal.File{Filename: file.Name()})
	}
	if err = templates.Index(entries).Render(r.Context(), w); err != nil {
		log.Println(err)
	}
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.TOC
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func sanitizeHTML(rawHtml []byte) []byte {
	p := bluemonday.UGCPolicy()
	return p.SanitizeBytes(rawHtml)
}

func renderFile(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")
	md, err := os.ReadFile("notes/" + filename)
	if err != nil {
		log.Println(err)
	}
	rawHtml := mdToHTML(md)
	sanitizedHtml := sanitizeHTML(rawHtml)
	if err = templates.File(string(sanitizedHtml)).Render(r.Context(), w); err != nil {
		log.Println(err)
	}
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /notes", index)
	mux.HandleFunc("GET /notes/{filename}", renderFile)

	fs := http.FileServer(http.Dir("assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fs))

	return mux
}

func main() {
	server := http.Server{
		Addr:    ":8000",
		Handler: logging(newMux()),
	}
	log.Fatal(server.ListenAndServe())
}
