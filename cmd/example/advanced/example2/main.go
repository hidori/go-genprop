package main

import (
	"fmt"
	"log"

	"github.com/hidori/go-genprop/example/advanced/example2"
)

//go:generate go run ../../../../genprop/main.go ../../../../example/advanced/example2/example.go > ../../../../example/advanced/example2/example_generated.go

func main() {
	fmt.Println("=== Advanced Example 2: Validation Demo ===")

	// Create a new user with valid data
	user, err := example2.NewUser("valid@example.com", "validpassword123", 85)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	// Display initial values
	fmt.Printf("User Email: %s\n", user.GetEmail())
	fmt.Printf("User Score: %d\n", user.GetScore())

	// Test valid updates
	fmt.Println("\n--- Testing valid updates ---")

	err = user.SetEmail("newemail@example.com")
	if err != nil {
		fmt.Printf("❌ Email update failed: %v\n", err)
	} else {
		fmt.Printf("✅ Email updated to: %s\n", user.GetEmail())
	}

	err = user.SetScore(95)
	if err != nil {
		fmt.Printf("❌ Score update failed: %v\n", err)
	} else {
		fmt.Printf("✅ Score updated to: %d\n", user.GetScore())
	}

	// Test invalid updates to demonstrate validation
	fmt.Println("\n--- Testing validation (these should fail) ---")

	err = user.SetEmail("invalid-email")
	if err != nil {
		fmt.Printf("✅ Expected validation error for invalid email: %v\n", err)
	}

	err = user.SetScore(150) // Out of range
	if err != nil {
		fmt.Printf("✅ Expected validation error for invalid score: %v\n", err)
	}

	// Final state
	fmt.Printf("\nFinal state:\n")
	fmt.Printf("User Email: %s\n", user.GetEmail())
	fmt.Printf("User Score: %d\n", user.GetScore())

	fmt.Println("\n✅ Advanced example 2 with validation completed successfully!")
}
