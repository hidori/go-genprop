package data

import (
	"fmt"
	"go/ast"
	"regexp"
)

type SuccessStruct struct {
	ignored1 int
	ignored2 int `prop:""`
	ignored3 int `prop:"-"`

	int1 int  `prop:"get,set"`
	int2 *int `prop:"get,set"`

	string1 string  `prop:"get,set"`
	string2 *string `prop:"get,set"`

	interface1 interface{}  `prop:"get,set"`
	interface2 *interface{} `prop:"get,set"`

	otherSuccessStruct1 OtherSuccessStruct  `prop:"get,set"`
	otherSuccessStruct2 *OtherSuccessStruct `prop:"get,set"`

	astFile1 ast.File  `prop:"get,set"`
	astFile2 *ast.File `prop:"get,set"`

	api         string `prop:"get"`
	apiEndpoint string `prop:"get"`
}

type OtherSuccessStruct struct {
	otherInt1 int `prop:"get"`
	otherInt2 int `prop:"get"`
}

type EmptyStruct struct{}

type String string

func Func() {
	type IgnoredInnerStruct struct {
		Int1 int `prop:"get"`
	}

	fmt.Println("Hello, world")
}

const IgnoredConst = 1

var IgnoredVar = regexp.MustCompile(`^$`)

var IgnoredAnonymousStruct = struct {
	Int1 int `prop:"get"`
}{
	Int1: 1,
}
