package main

import (
	"log"
	"os"

	"github.com/monochromegane/dragon-imports"
)

func main() {
	err := dragon.Imports()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
