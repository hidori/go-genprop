package main

import (
	"fmt"
	"net/http"

	example "github.com/hidori/go-genprop/example/initialism"
)

func main() {
	// Create HTTP client
	httpClient := http.Client{}

	// Create a new API client using the constructor
	client := example.NewAPIClient(httpClient, "https://api.example.com", "secret-api-key")

	// Access values using getters
	fmt.Printf("HTTP Client: %+v\n", client.GetClient()) // Output: HTTP Client: {Transport:<nil> CheckRedirect:<nil> Jar:<nil> Timeout:0s}
	fmt.Printf("Client URL: %s\n", client.GetURL())      // Output: Client URL: https://api.example.com
	fmt.Printf("API Key: %s\n", client.GetAPIKey())      // Output: API Key: secret-api-key

	// Note: setters are private, so they can only be called from within the package
	// This demonstrates initialism handling with private setters for controlled access
}
