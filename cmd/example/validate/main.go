package main

import (
	"fmt"
	"log"

	example "github.com/hidori/go-genprop/example/validate"
)

func main() {
	// Create a new account using the generated New function with validation
	account, err := example.NewAccount(1, "john_doe", "john@example.com", "SecurePass123!")
	if err != nil {
		log.Fatal("Failed to create account:", err)
	}

	// Access values using getters
	fmt.Printf("Account ID: %d\n", account.GetID())     // Output: Account ID: 1
	fmt.Printf("Username: %s\n", account.GetUsername()) // Output: Username: john_doe
	fmt.Printf("Email: %s\n", account.GetEmail())       // Output: Email: john@example.com

	// Update values using setters with custom validation
	err = account.SetUsername("jane_smith")
	if err != nil {
		fmt.Printf("Username update failed: %v\n", err)
	} else {
		fmt.Printf("Updated Username: %s\n", account.GetUsername()) // Output: Updated Username: jane_smith
	}

	err = account.SetEmail("jane@newdomain.com")
	if err != nil {
		fmt.Printf("Email update failed: %v\n", err)
	} else {
		fmt.Printf("Updated Email: %s\n", account.GetEmail()) // Output: Updated Email: jane@newdomain.com
	}

	// Test validation failures
	fmt.Println("\nTesting validation failures:")

	err = account.SetUsername("ab") // Too short
	if err != nil {
		fmt.Printf("Expected validation error: %v\n", err)
	}

	err = account.SetEmail("invalid-email") // Invalid format
	if err != nil {
		fmt.Printf("Expected validation error: %v\n", err)
	}

	// Note: password setter is private and uses validation
	fmt.Println("\nNote: Password setter is private and validates strong password requirements")
}
