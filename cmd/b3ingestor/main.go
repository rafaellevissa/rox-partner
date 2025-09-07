package main

import (
	"log"

	"github.com/rafaellevissa/rox-partner/internal/cli"
)

func main() {
	root := cli.NewRootCmd()
	if err := root.Execute(); err != nil {
		log.Fatalf("command error: %v", err)
	}
}
