package main

import (
	"fmt"

	"github.com/hidori/go-genprop/example/basic"
)

//go:generate go run ../../../genprop/main.go ../../../example/basic/example.go > ../../../example/basic/example_generated.go

func main() {
	fmt.Println("=== Basic Example Demo ===")

	// Create a new user using the constructor
	user := basic.NewUser(1, "John Doe", "john@example.com", "secretpassword")

	// Display initial values using getters
	fmt.Printf("User ID: %d\n", user.GetID())
	fmt.Printf("User Name: %s\n", user.GetName())
	fmt.Printf("User Email: %s\n", user.GetEmail())

	// Update values using setters
	user.SetName("Jane Smith")
	user.SetEmail("jane@example.com")

	// Display updated values
	fmt.Printf("\nAfter updates:\n")
	fmt.Printf("User Name: %s\n", user.GetName())
	fmt.Printf("User Email: %s\n", user.GetEmail())

	fmt.Println("\n✅ Basic example completed successfully!")
}
