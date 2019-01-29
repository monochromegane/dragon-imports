package main

import (
	"flag"
	"log"
	"os"

	dragon "github.com/monochromegane/dragon-imports"
)

var restore bool

func init() {
	flag.BoolVar(&restore, "restore", false, "goimports is returned to original one.")
}

func main() {
	flag.Parse()
	err := dragon.Imports(restore)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
