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

func ExampleValidation() {
	v1 := Valid[int, string](10)
	v2 := Invalid[int, string]("error 1")
	v3 := Invalid[int, string]("error 2")

	res := v1.And(v2, func(a, b int) int { return a + b }).
		And(v3, func(a, b int) int { return a + b })

	fmt.Println(res.Errors())
	// Output: [error 1 error 2]
}

func ExampleReader() {
	type Env struct{ Port int }
	r := Ask[Env]().Map(func(e Env) string {
		return fmt.Sprintf("Port: %d", e.Port)
	})

	fmt.Println(r.Run(Env{Port: 8080}))
	// Output: Port: 8080
}

func ExampleWriter() {
	w := NewWriter(10, "start").
		AndThen(func(i int) Writer[string, int] {
			return NewWriter(i*2, "doubled")
		})

	fmt.Printf("Value: %d, Logs: %v\n", w.Value(), w.Logs())
	// Output: Value: 20, Logs: [start doubled]
}

func ExampleState() {
	program := Modify(func(s int) int { return s + 1 }).
		AndThen(func(_ any) State[int, int] {
			return Get[int]()
		})

	val, _ := program.Run(10)
	fmt.Println(val)
	// Output: 11
}

func ExampleTask() {
	t := NewTask(func() Result[int] {
		return Ok(42)
	}).Map(func(i int) int { return i * 2 })

	res := t.Run()
	val, _ := res.Get()
	fmt.Println(val)
	// Output: 84
}
