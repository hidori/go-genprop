package main

import (
	"fmt"

	example "github.com/hidori/go-genprop/example/basic"
)

func main() {
	// Create a new user using the constructor
	user := example.NewUser(1, "John Doe")

	// Access values using getters
	fmt.Printf("User ID: %d\n", user.GetID())     // Output: User ID: 1
	fmt.Printf("User Name: %s\n", user.GetName()) // Output: User Name: John Doe

	// Update values using setters
	user.SetName("Jane Smith")
	fmt.Printf("Updated Name: %s\n", user.GetName()) // Output: Updated Name: Jane Smith
}
