package main

import (
	"fmt"
	"log"

	"github.com/hidori/go-genprop/example/advanced"
)

//go:generate go run ../../../genprop/main.go ../../../example/advanced/user.go > ../../../example/advanced/user_prop.go

func main() {
	// Create a new user with validation
	user, err := advanced.NewUser(1, "John Doe", "john@example.com")
	if err != nil {
		log.Fatal(err)
	}

	// Access values using getters
	fmt.Printf("User ID: %d\n", user.GetID())       // Output: User ID: 1
	fmt.Printf("User Name: %s\n", user.GetName())   // Output: User Name: John Doe
	fmt.Printf("User Email: %s\n", user.GetEmail()) // Output: User Email: john@example.com

	// Update values using setters with validation
	user.SetName("Jane Smith")

	// Note: setEmail is private, so it can only be called from within the package
	// This demonstrates validation with private setters for controlled access
}
