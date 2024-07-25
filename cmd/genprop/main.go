package main

import (
	"log"

	app "github.com/hidori/go-genprop/app/genprop"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
