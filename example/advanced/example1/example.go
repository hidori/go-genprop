package example1

type Person struct {
	id   int    `property:"get"`
	name string `property:"get,set"`
	age  int    `property:"get,set"`
}

func NewPerson(id int, name string, age int) *Person {
	person := &Person{id: id}

	person.SetName(name)
	person.SetAge(age)

	return person
}
