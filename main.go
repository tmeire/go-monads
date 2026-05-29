package main

import "fmt"

func main() {
	e := Success[int, error](42).
		AndThen(func(i int) Expected[int, error] {
			return Success[int, error](i + 1)
		}).
		AndThen(func(i int) Expected[int, error] {
			return Success[int, error](i * 2)
		}).
		Transform(func(i int) string {
			return fmt.Sprintf("-%d-", i)
		}).
		AndThen(func(s string) Expected[int, error] {
			return Failure[int, error](fmt.Errorf("andThen error"))
		}).
		//OrElse(func(e error) Expected[int, error] {
		//	return Success[int, error](67)
		//}).
		TransformError(func(e error) string {
			return fmt.Errorf("transformed error: %w", e).Error()
		})

	fmt.Printf("%#v\n", e)
}
