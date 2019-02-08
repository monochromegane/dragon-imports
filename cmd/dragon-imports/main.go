package main

import (
	"flag"
	"log"
	"os"

	dragon "github.com/monochromegane/dragon-imports"
)

var (
	restore   bool
	gomodules bool
)

func init() {
	flag.BoolVar(&restore, "restore", false, "goimports is returned to original one.")
	flag.BoolVar(&gomodules, "gomodules", false, "goimports makes zstdlib.go from Go modules packages instead of GOPATH/src.")
}

func main() {
	flag.Parse()
	err := dragon.Imports(restore, gomodules)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
