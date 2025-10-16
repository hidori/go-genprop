package example

// User represents a user with basic information.
type User struct {
	id    int    `property:"get"`     // Read-only ID field
	name  string `property:"get,set"` // Name with both getter and setter
	email string `property:"get,set"` // Email with both getter and setter
}
