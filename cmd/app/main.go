package main

import (
	"log"

	"github.com/annguyen17-tiki/grab/internal/handler"
)

func main() {
	server, err := handler.New()
	if err != nil {
		log.Fatalf("failed to prepare handler, err: %v", err)
	}

	err = server.Run()
	if err != nil {
		log.Fatalf("failed to run server, err: %v", err)
	}
}
