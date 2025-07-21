package main

import (
	"fmt"

	example "github.com/hidori/go-genprop/example/new"
)

func main() {
	// Create a new user using the generated New function
	user := example.NewUser(1, "John Doe", "john@example.com")

	// Access values using getters
	fmt.Printf("User ID: %d\n", user.GetID())       // Output: User ID: 1
	fmt.Printf("User Name: %s\n", user.GetName())   // Output: User Name: John Doe
	fmt.Printf("User Email: %s\n", user.GetEmail()) // Output: User Email: john@example.com

	// Update values using setters
	user.SetName("Jane Smith")
	user.SetEmail("jane@example.com")

	// Display updated values
	fmt.Printf("Updated Name: %s\n", user.GetName())   // Output: Updated Name: Jane Smith
	fmt.Printf("Updated Email: %s\n", user.GetEmail()) // Output: Updated Email: jane@example.com
}
