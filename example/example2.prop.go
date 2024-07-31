// Code generated by github.com/hidori/go-genprop/cmd/genprop DO NOT EDIT.
package example

func (t *Example2Struct) GetValue1() int {
	return t.value1
}
func (t *Example2Struct) SetValue2(v int) {
	t.value2 = v
}
func (t *Example2Struct) GetValue3() int {
	return t.value3
}
func (t *Example2Struct) SetValue3(v int) {
	t.value3 = v
}
func (t *Example2Struct) GetValue4() int {
	return t.value4
}
func (t *Example2Struct) SetValue4(v int) error {
	err := validateFieldValue(v, "required,min=1")
	if err != nil {
		return err
	}
	t.value4 = v
	return nil
}
func (t *Example2Struct) GetID() int {
	return t.id
}
func (t *Example2Struct) GetAPI() string {
	return t.api
}
func (t *Example2Struct) GetURL() string {
	return t.url
}
func (t *Example2Struct) GetHttp() string {
	return t.http
}