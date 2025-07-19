package main

import (
	"log"

	"github.com/hidori/go-genprop/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatalf("Error running generator: %v", err)
	}
}
