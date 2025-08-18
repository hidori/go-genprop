package example

// User represents a basic user with ID and name.
type User struct {
	id   int    `property:"get"`     // Read-only ID field
	name string `property:"get,set"` // Name with both getter and setter
}

// NewUser creates a new User instance
func NewUser(id int, name string) *User {
	user := &User{id: id}
	user.SetName(name)

	return user
}
