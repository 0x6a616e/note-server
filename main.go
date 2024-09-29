package main

import (
	"log"
	"net/http"
	"os"
	"strings"
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

func renderFolder(w http.ResponseWriter, r *http.Request, folder string) {
	if !strings.HasSuffix(folder, "/") {
		folder += "/"
	}
	files, err := os.ReadDir(folder)
	if err != nil {
		log.Println(err)
	}
	entries := []internal.File{}
	for _, file := range files {
		filename := folder + file.Name()
		entries = append(entries, internal.File{Filename: filename})
	}
	foldername := internal.File{Filename: folder[:len(folder)-1]}
	if err = templates.Index(foldername, entries).Render(r.Context(), w); err != nil {
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

func renderFile(w http.ResponseWriter, r *http.Request, filename string) {
	md, err := os.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}
	rawHtml := mdToHTML(md)
	sanitizedHtml := sanitizeHTML(rawHtml)
	if err = templates.File(string(sanitizedHtml)).Render(r.Context(), w); err != nil {
		log.Println(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Println(err)
	}
	if fileInfo.IsDir() {
		renderFolder(w, r, filename)
	} else {
		renderFile(w, r, filename)
	}
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fs))

	mux.HandleFunc("GET /files/{filename...}", index)

	return mux
}

func main() {
	server := http.Server{
		Addr:    ":8000",
		Handler: logging(newMux()),
	}
	log.Fatal(server.ListenAndServe())
}
