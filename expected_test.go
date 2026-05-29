package main

import (
	"errors"
	"testing"
)

func TestAndThen(t *testing.T) {
	t.Run("Value to Value", func(t *testing.T) {
		e := Success[int, string](42)
		res := e.AndThen(func(i int) Expected[string, string] {
			return Success[string, string]("success")
		})
		if res.t != "success" || !res.hasValue {
			t.Errorf("expected success value, got %#v", res)
		}
	})

	t.Run("Value to Error", func(t *testing.T) {
		e := Success[int, string](42)
		res := e.AndThen(func(i int) Expected[string, string] {
			return Failure[string, string]("andthen failure")
		})
		if res.e != "andthen failure" || res.hasValue {
			t.Errorf("expected failure, got %#v", res)
		}
	})

	t.Run("Error remains Error", func(t *testing.T) {
		e := Failure[int, string]("initial failure")
		res := e.AndThen(func(i int) Expected[string, string] {
			return Success[string, string]("should not happen")
		})
		if res.e != "initial failure" || res.hasValue {
			t.Errorf("expected initial failure, got %#v", res)
		}
	})
}

func TestTransform(t *testing.T) {
	t.Run("Value", func(t *testing.T) {
		e := Success[int, string](42)
		res := e.Transform(func(i int) string {
			return "transformed"
		})
		if res.t != "transformed" || !res.hasValue {
			t.Errorf("expected transformed value, got %#v", res)
		}
	})

	t.Run("Error", func(t *testing.T) {
		e := Failure[int, string]("failure")
		res := e.Transform(func(i int) string {
			return "should not happen"
		})
		if res.e != "failure" || res.hasValue {
			t.Errorf("expected failure, got %#v", res)
		}
	})
}

func TestOrElse(t *testing.T) {
	t.Run("Value remains Value", func(t *testing.T) {
		e := Success[int, string](42)
		res := e.OrElse(func(s string) Expected[int, int] {
			return Success[int, int](100)
		})
		if res.t != 42 || !res.hasValue {
			t.Errorf("expected 42, got %#v", res)
		}
	})

	t.Run("Error to Value", func(t *testing.T) {
		e := Failure[int, string]("failure")
		res := e.OrElse(func(s string) Expected[int, int] {
			return Success[int, int](100)
		})
		if res.t != 100 || !res.hasValue {
			t.Errorf("expected 100, got %#v", res)
		}
	})

	t.Run("Error to Error", func(t *testing.T) {
		e := Failure[int, string]("initial failure")
		res := e.OrElse(func(s string) Expected[int, int] {
			return Failure[int, int](500)
		})
		if res.e != 500 || res.hasValue {
			t.Errorf("expected 500, got %#v", res)
		}
	})
}

func TestTransformError(t *testing.T) {
	t.Run("Value", func(t *testing.T) {
		e := Success[int, string](42)
		res := e.TransformError(func(s string) int {
			return 500
		})
		if res.t != 42 || !res.hasValue {
			t.Errorf("expected 42, got %#v", res)
		}
	})

	t.Run("Error", func(t *testing.T) {
		e := Failure[int, string]("failure")
		res := e.TransformError(func(s string) int {
			return 500
		})
		if res.e != 500 || res.hasValue {
			t.Errorf("expected 500, got %#v", res)
		}
	})
}

func TestWithErrorInterface(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")

	t.Run("AndThen", func(t *testing.T) {
		e := Success[int, error](42)
		res := e.AndThen(func(i int) Expected[string, error] {
			return Failure[string, error](err1)
		})
		if res.e != err1 {
			t.Errorf("expected err1, got %v", res.e)
		}
	})

	t.Run("OrElse", func(t *testing.T) {
		e := Failure[int, error](err1)
		res := e.OrElse(func(err error) Expected[int, error] {
			return Failure[int, error](err2)
		})
		if res.e != err2 {
			t.Errorf("expected err2, got %v", res.e)
		}
	})
}

func TestZeroValuesFixed(t *testing.T) {
	t.Run("Int Zero Value", func(t *testing.T) {
		// e is value 0
		e := Success[int, string](0)
		
		res := e.OrElse(func(s string) Expected[int, string] {
			return Success[int, string](100)
		})
		
		if res.t != 0 || !res.hasValue {
			t.Errorf("Expected 0 to be preserved, got %#v", res)
		}
	})

	t.Run("String Zero Value as Error", func(t *testing.T) {
		// e is error ""
		e := Failure[int, string]("")
		
		res := e.AndThen(func(i int) Expected[int, string] {
			return Success[int, string](i + 1)
		})
		
		if res.hasValue {
			t.Errorf("Empty string error should have blocked AndThen, but it proceeded")
		}
		if res.e != "" {
			t.Errorf("Expected error '', got %v", res.e)
		}
	})
}

func TestNonComparable(t *testing.T) {
	t.Run("Slice as Value", func(t *testing.T) {
		e := Success[[]int, string]([]int{1, 2, 3})
		res := e.Transform(func(s []int) int {
			return len(s)
		})
		if res.t != 3 || !res.hasValue {
			t.Errorf("Expected length 3, got %#v", res)
		}
	})
}
