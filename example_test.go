package main

import (
	"fmt"
	"strings"
)

func ExampleOptional() {
	o := Some(42)
	res := o.Filter(func(i int) bool { return i > 40 }).
		Transform(func(i int) string { return fmt.Sprintf("Result: %d", i) })

	fmt.Println(res.ValueOr("Default"))
	// Output: Result: 42
}

func ExampleExpected() {
	e := Success[int, string](42)
	res := e.AndThen(func(i int) Expected[int, string] {
		if i > 50 {
			return Failure[int, string]("too large")
		}
		return Success[int, string](i * 2)
	})

	res.IfHasValueOrElse(
		func(i int) { fmt.Printf("Value: %d\n", i) },
		func(s string) { fmt.Printf("Error: %s\n", s) },
	)
	// Output: Value: 84
}

func ExampleResult() {
	r := From(42, nil)
	val := r.Map(func(i int) string {
		return strings.Repeat("*", i/10)
	}).ValueOr("failed")

	fmt.Println(val)
	// Output: ****
}

func ExampleEither() {
	e := Right[int, string]("hello")
	val := e.Match(
		func(i int) string { return fmt.Sprintf("Int: %d", i) },
		func(s string) string { return "String: " + s },
	)

	fmt.Println(val)
	// Output: String: hello
}

func ExampleTry() {
	t := Invoke(func() (int, error) {
		return 42, nil
	}).Map(func(i int) int {
		if i == 42 {
			panic("it's the answer")
		}
		return i
	})

	_, err := t.Get()
	fmt.Println(err)
	// Output: panic: it's the answer
}

func ExampleLazy() {
	l := NewLazy(func() int {
		fmt.Println("Initializing...")
		return 42
	})

	fmt.Println("First access:")
	fmt.Println(l.Get())
	fmt.Println("Second access:")
	fmt.Println(l.Get())
	// Output:
	// First access:
	// Initializing...
	// 42
	// Second access:
	// 42
}

func ExampleStream() {
	s := OfSlice([]int{1, 2, 3, 4, 5})
	res := s.Filter(func(i int) bool { return i%2 != 0 }).
		Map(func(i int) int { return i * i }).
		ToSlice()

	fmt.Println(res)
	// Output: [1 9 25]
}
