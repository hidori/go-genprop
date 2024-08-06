package main

import (
	"fmt"

	"github.com/hidori/go-genprop/example"
)

func main() {
	v1, err := example.NewStruct(1, 2, 3, 4, 5)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(v1)
	}

	v2, err := example.NewStruct(1, 2, 0, 4, -5)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(v2)
	}
}
