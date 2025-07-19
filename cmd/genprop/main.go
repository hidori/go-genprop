package main

import (
	"log"
	"os"

	"github.com/hidori/go-genprop/internal/app"
)

func main() {
	err := app.Run(os.Args[1:])
	if err != nil {
		log.Fatalf("Error running generator: %v", err)
	}
}
