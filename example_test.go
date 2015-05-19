package pgparser

import (
	"fmt"
)

func ExampleParse() {
	var s = `(a,{x,y,z},5)`

	var v struct {
		A string
		B []string
		C int
	}

	if err := Unmarshal(s, &v); err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", v)

	// Output:
	// struct { A string; B []string; C int }{A:"a", B:[]string{"x", "y", "z"}, C:5}
}
