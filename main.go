package main

import "fmt"

type Test struct {
	Name string
}

func (t *Test) Print[T fmt.Stringer](val T) {
	println(t.Name, val.String())
}

type S string

func (s S) String() string {
	return string(s)
}

func main() {
	println("Hello, World!")

	t := Test{Name: "Name"}
	t.Print(S("val"))
}
