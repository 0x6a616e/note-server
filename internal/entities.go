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
	filename := strings.TrimSuffix(f.Filename, ".md")
	words := strings.Split(filename, "_")
	filename = strings.Join(words, " ")
	caser := cases.Title(language.English)
	filename = caser.String(filename)
	return filename
}
