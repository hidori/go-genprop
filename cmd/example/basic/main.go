package main

import (
	"fmt"

	"github.com/hidori/go-genprop/example/basic"
)

//go:generate go run ../../../genprop/main.go ../../../example/basic/user.go > ../../../example/basic/user_prop.go

func main() {
	// Create a new user using the constructor
	user := basic.NewUser(1, "John Doe")

	// Access values using getters
	fmt.Printf("User ID: %d\n", user.GetID())     // Output: User ID: 1
	fmt.Printf("User Name: %s\n", user.GetName()) // Output: User Name: John Doe

	// Update values using setters
	user.SetName("Jane Smith")
	fmt.Printf("Updated Name: %s\n", user.GetName()) // Output: Updated Name: Jane Smith
}
