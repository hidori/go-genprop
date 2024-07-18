package data

import (
	"fmt"
	"go/ast"
	"regexp"
)

type SuccessStruct struct {
	ignored1 int
	ignored2 int `property:""`
	ignored3 int `property:"-"`

	int1 int  `property:"get,set"`
	int2 *int `property:"get,set"`

	string1 string  `property:"get,set"`
	string2 *string `property:"get,set"`

	interface1 interface{}  `property:"get,set"`
	interface2 *interface{} `property:"get,set"`

	otherSuccessStruct1 OtherSuccessStruct  `property:"get,set"`
	otherSuccessStruct2 *OtherSuccessStruct `property:"get,set"`

	astFile1 ast.File  `property:"get,set"`
	astFile2 *ast.File `property:"get,set"`

	api         string `property:"get"`
	apiEndpoint string `property:"get"`
}

type OtherSuccessStruct struct {
	otherInt1 int `property:"get"`
	otherInt2 int `property:"get"`
}

type EmptyStruct struct{}

type String string

func Func() {
	type IgnoredInnerStruct struct {
		Int1 int `property:"get"`
	}

	fmt.Println("Hello, world")
}

const IgnoredConst = 1

var IgnoredVar = regexp.MustCompile(`^$`)

var IgnoredAnonymousStruct = struct {
	Int1 int `property:"get"`
}{
	Int1: 1,
}
