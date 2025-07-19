package main

import (
	"fmt"

	"github.com/hidori/go-genprop/example/advanced/example1"
)

//go:generate go run ../../../../genprop/main.go ../../../../example/advanced/example1/example.go > ../../../../example/advanced/example1/example_generated.go

func main() {
	fmt.Println("=== Advanced Example 1: Basic Structure Demo ===")

	// Create a new person using the constructor
	person := example1.NewPerson(100, "Alice Johnson", 30)

	// Display initial values using getters
	fmt.Printf("Person ID: %d\n", person.GetID())
	fmt.Printf("Person Name: %s\n", person.GetName())
	fmt.Printf("Person Age: %d\n", person.GetAge())

	// Update values using setters
	person.SetName("Alice Williams")
	person.SetAge(31)

	// Display updated values
	fmt.Printf("\nAfter updates:\n")
	fmt.Printf("Person Name: %s\n", person.GetName())
	fmt.Printf("Person Age: %d\n", person.GetAge())

	// ID is read-only
	fmt.Printf("Person ID (read-only): %d\n", person.GetID())

	fmt.Println("\n✅ Advanced example 1 completed successfully!")
}
