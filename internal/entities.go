package internal

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type File struct {
	Filename string
}

func (f File) String() string {
	pathParts := strings.Split(f.Filename, "/")
	filename := pathParts[len(pathParts)-1]
	filename = strings.TrimSuffix(filename, ".md")
	words := strings.Split(filename, "_")
	filename = strings.Join(words, " ")
	caser := cases.Title(language.English)
	filename = caser.String(filename)
	return filename
}
