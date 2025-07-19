package basic

type User struct {
	id       int    `property:"get"`
	name     string `property:"get,set"`
	email    string `property:"get,set"`
	password string `property:"set=private"`
}

// NewUser creates a new User instance
func NewUser(id int, name, email, password string) *User {
	user := &User{id: id}

	user.SetName(name)
	user.SetEmail(email)
	user.setPassword(password)

	return user
}
