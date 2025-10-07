package example

import "net/http"

// APIClient represents an API client with various initialism fields.
type APIClient struct {
	client http.Client `property:"get,set=private"` // HTTP client
	url    string      `property:"get,set=private"` // URL field (initialism)
}

// NewAPIClient creates a new APIClient instance
func NewAPIClient(client http.Client, url string) *APIClient {
	apiClient := &APIClient{}
	apiClient.setClient(client)
	apiClient.setURL(url)

	return apiClient
}
