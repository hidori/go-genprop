// Code generated by github.com/hidori/go-genprop/cmd/genprop DO NOT EDIT.
package advanced

func (t *User) GetID() int {
	return t.id
}
func (t *User) GetName() string {
	return t.name
}
func (t *User) SetName(v string) {
	t.name = v
}
func (t *User) GetEmail() string {
	return t.email
}
func (t *User) setEmail(v string) error {
	err := validateFieldValue("email", v, "required,email")
	if err != nil {
		return err
	}
	t.email = v
	return nil
}
