package main

import (
	"embed"
	"fmt"
	"os"

	new "github.com/shindakun/goproj/cmd"
)

//go:embed templates/*.tmpl
var tpls embed.FS

func usage() {
	fmt.Printf(`
Goproj will create a new Go project repo.

Usage:

	goproj new <directory name>
	`)
}

func main() {

	if len(os.Args) < 3 {
		usage()
		os.Exit(0)
	}
	switch cmd := os.Args[1]; cmd {
	case "new":
		dir := os.Args[2]
		new.CmdNew(dir, tpls)
	case "help":
		usage()
		os.Exit(0)
	default:
		usage()
		os.Exit(0)
	}
}
