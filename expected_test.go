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

func TestObservers(t *testing.T) {
	t.Run("HasValue", func(t *testing.T) {
		if !Success[int, string](42).HasValue() {
			t.Error("expected Success to have value")
		}
		if Failure[int, string]("err").HasValue() {
			t.Error("expected Failure not to have value")
		}
	})

	t.Run("ValueOr", func(t *testing.T) {
		if Success[int, string](42).ValueOr(100) != 42 {
			t.Error("expected 42")
		}
		if Failure[int, string]("err").ValueOr(100) != 100 {
			t.Error("expected 100")
		}
	})

	t.Run("ErrorOr", func(t *testing.T) {
		if Failure[int, string]("err").ErrorOr("default") != "err" {
			t.Error("expected 'err'")
		}
		if Success[int, string](42).ErrorOr("default") != "default" {
			t.Error("expected 'default'")
		}
	})

	t.Run("Get", func(t *testing.T) {
		val, err := Success[int, string](42).Get()
		if val != 42 || err != "" {
			t.Errorf("expected (42, ''), got (%d, '%s')", val, err)
		}

		val, err = Failure[int, string]("err").Get()
		if val != 0 || err != "err" {
			t.Errorf("expected (0, 'err'), got (%d, '%s')", val, err)
		}
	})
}

func TestEmplace(t *testing.T) {
	t.Run("Emplace on Success", func(t *testing.T) {
		e := Success[int, string](42)
		valPtr := e.Emplace(100)
		if !e.hasValue || e.t != 100 || e.e != "" {
			t.Errorf("expected success with 100, got %#v", e)
		}
		if *valPtr != 100 {
			t.Errorf("expected returned pointer to point to 100, got %d", *valPtr)
		}
		// Modify via pointer
		*valPtr = 200
		if e.t != 200 {
			t.Errorf("expected internal value to be 200 after pointer modification, got %d", e.t)
		}
	})

	t.Run("Emplace on Failure", func(t *testing.T) {
		e := Failure[int, string]("error")
		e.Emplace(42)
		if !e.hasValue || e.t != 42 || e.e != "" {
			t.Errorf("expected success with 42, got %#v", e)
		}
	})

	t.Run("EmplaceError on Success", func(t *testing.T) {
		e := Success[int, string](42)
		errPtr := e.EmplaceError("new error")
		if e.hasValue || e.e != "new error" || e.t != 0 {
			t.Errorf("expected failure with 'new error', got %#v", e)
		}
		if *errPtr != "new error" {
			t.Errorf("expected returned pointer to point to 'new error', got %s", *errPtr)
		}
	})

	t.Run("EmplaceError on Failure", func(t *testing.T) {
		e := Failure[int, string]("old error")
		e.EmplaceError("new error")
		if e.hasValue || e.e != "new error" || e.t != 0 {
			t.Errorf("expected failure with 'new error', got %#v", e)
		}
	})
}

func TestExpectedUtilities(t *testing.T) {
	t.Run("ValueOrGet", func(t *testing.T) {
		if Success[int, string](42).ValueOrGet(func() int { return 100 }) != 42 {
			t.Error("expected 42")
		}
		if Failure[int, string]("err").ValueOrGet(func() int { return 100 }) != 100 {
			t.Error("expected 100")
		}
	})

	t.Run("IfHasValue", func(t *testing.T) {
		called := false
		Success[int, string](42).IfHasValue(func(i int) {
			called = true
		})
		if !called {
			t.Error("expected function to be called")
		}

		called = false
		Failure[int, string]("err").IfHasValue(func(i int) {
			called = true
		})
		if called {
			t.Error("expected function not to be called")
		}
	})

	t.Run("IfHasError", func(t *testing.T) {
		called := false
		Failure[int, string]("err").IfHasError(func(s string) {
			called = true
		})
		if !called {
			t.Error("expected function to be called")
		}

		called = false
		Success[int, string](42).IfHasError(func(s string) {
			called = true
		})
		if called {
			t.Error("expected function not to be called")
		}
	})

	t.Run("IfHasValueOrElse", func(t *testing.T) {
		state := ""
		Success[int, string](42).IfHasValueOrElse(
			func(i int) { state = "value" },
			func(s string) { state = "error" },
		)
		if state != "value" {
			t.Error("expected 'value'")
		}

		Failure[int, string]("err").IfHasValueOrElse(
			func(i int) { state = "value" },
			func(s string) { state = "error" },
		)
		if state != "error" {
			t.Error("expected 'error'")
		}
	})

	t.Run("Equals", func(t *testing.T) {
		e1 := Success[[]int, string]([]int{1, 2})
		e2 := Success[[]int, string]([]int{1, 2})
		e3 := Success[[]int, string]([]int{1, 3})
		e4 := Failure[[]int, string]("err")
		e5 := Failure[[]int, string]("err")
		e6 := Failure[[]int, string]("other")

		if !e1.Equals(e2) {
			t.Error("e1 should equal e2")
		}
		if e1.Equals(e3) {
			t.Error("e1 should not equal e3")
		}
		if e1.Equals(e4) {
			t.Error("e1 should not equal e4")
		}
		if !e4.Equals(e5) {
			t.Error("e4 should equal e5")
		}
		if e4.Equals(e6) {
			t.Error("e4 should not equal e6")
		}
	})
}
