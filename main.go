package main

import (
	"embed"
	"os"

	new "github.com/shindakun/goproj/cmd"
)

//go:embed templates/*.tmpl
var tpls embed.FS

func main() {

	// goproj new {{directory}}

	if len(os.Args) < 3 {
		panic("usage instructions")
	}
	switch cmd := os.Args[1]; cmd {
	case "new":
		dir := os.Args[2]
		new.CmdNew(dir, tpls)
	case "help":
		panic("usage instructions")
	default:
		panic("usage instructions")
	}
}
