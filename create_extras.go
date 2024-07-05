package main

import (
	_ "embed"
	"io"
	"os"
	"strings"
)

//go:embed extras/.gitignore
var gitignore string

//go:embed extras/README.md
var README string

func GenerateExtra(file, content string) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	io.Copy(f, strings.NewReader(content))
}

func GenerateExtras() {
	GenerateExtra(".gitignore", gitignore)
	GenerateExtra("README.md", README)
}
