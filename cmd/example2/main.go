package main

import (
	"fmt"

	"github.com/hidori/go-genprop/example"
)

func main() {
	v1, err := example.NewExample2Struct(1, 2, 3)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(v1)
	}

	v2, err := example.NewExample2Struct(1, 2, -3)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(v2)
	}
}
