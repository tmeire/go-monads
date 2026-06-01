package main

import (
	"errors"
	"testing"
)

func TestOptional(t *testing.T) {
	t.Run("Some and None", func(t *testing.T) {
		o1 := Some(42)
		if !o1.HasValue() || o1.value != 42 {
			t.Errorf("expected Some(42), got %#v", o1)
		}

		o2 := None[int]()
		if o2.HasValue() {
			t.Errorf("expected None, got %#v", o2)
		}
	})

	t.Run("ValueOr", func(t *testing.T) {
		if Some(42).ValueOr(100) != 42 {
			t.Error("expected 42")
		}
		if None[int]().ValueOr(100) != 100 {
			t.Error("expected 100")
		}
	})

	t.Run("Get", func(t *testing.T) {
		val, ok := Some(42).Get()
		if !ok || val != 42 {
			t.Errorf("expected (42, true), got (%d, %v)", val, ok)
		}

		val, ok = None[int]().Get()
		if ok || val != 0 {
			t.Errorf("expected (0, false), got (%d, %v)", val, ok)
		}
	})

	t.Run("Emplace", func(t *testing.T) {
		o := None[int]()
		valPtr := o.Emplace(42)
		if !o.HasValue() || o.value != 42 {
			t.Errorf("expected Some(42), got %#v", o)
		}
		if *valPtr != 42 {
			t.Errorf("expected returned pointer to point to 42, got %d", *valPtr)
		}

		*valPtr = 100
		if o.value != 100 {
			t.Errorf("expected internal value to be 100 after pointer modification, got %d", o.value)
		}

		o.Emplace(200)
		if o.value != 200 {
			t.Errorf("expected internal value to be 200, got %d", o.value)
		}
	})

	t.Run("Reset", func(t *testing.T) {
		o := Some(42)
		o.Reset()
		if o.HasValue() {
			t.Errorf("expected None after Reset, got %#v", o)
		}
	})

	t.Run("AndThen", func(t *testing.T) {
		o := Some(21)
		res := o.AndThen(func(i int) Optional[int] {
			return Some(i * 2)
		})
		if !res.HasValue() || res.value != 42 {
			t.Errorf("expected Some(42), got %#v", res)
		}

		res = None[int]().AndThen(func(i int) Optional[int] {
			return Some(i * 2)
		})
		if res.HasValue() {
			t.Error("expected None")
		}
	})

	t.Run("Transform", func(t *testing.T) {
		o := Some(42)
		res := o.Transform(func(i int) string {
			return "value"
		})
		if !res.HasValue() || res.value != "value" {
			t.Errorf("expected Some('value'), got %#v", res)
		}

		resNone := None[int]().Transform(func(i int) string {
			return "value"
		})
		if resNone.HasValue() {
			t.Error("expected None")
		}
	})

	t.Run("ValueOrGet", func(t *testing.T) {
		if Some(42).ValueOrGet(func() int { return 100 }) != 42 {
			t.Error("expected 42")
		}
		if None[int]().ValueOrGet(func() int { return 100 }) != 100 {
			t.Error("expected 100")
		}
	})

	t.Run("ValueOrErr", func(t *testing.T) {
		err := errors.New("missing")
		val, e := Some(42).ValueOrErr(err)
		if e != nil || val != 42 {
			t.Errorf("expected (42, nil), got (%d, %v)", val, e)
		}

		val, e = None[int]().ValueOrErr(err)
		if e != err || val != 0 {
			t.Errorf("expected (0, error), got (%d, %v)", val, e)
		}
	})

	t.Run("Filter", func(t *testing.T) {
		o := Some(42)
		res := o.Filter(func(i int) bool { return i > 40 })
		if !res.HasValue() || res.value != 42 {
			t.Error("expected Some(42)")
		}

		res = o.Filter(func(i int) bool { return i < 40 })
		if res.HasValue() {
			t.Error("expected None")
		}

		res = None[int]().Filter(func(i int) bool { return true })
		if res.HasValue() {
			t.Error("expected None")
		}
	})

	t.Run("IfPresent", func(t *testing.T) {
		called := false
		Some(42).IfPresent(func(i int) {
			called = true
		})
		if !called {
			t.Error("expected function to be called")
		}

		called = false
		None[int]().IfPresent(func(i int) {
			called = true
		})
		if called {
			t.Error("expected function not to be called")
		}
	})

	t.Run("IfPresentOrElse", func(t *testing.T) {
		called := ""
		Some(42).IfPresentOrElse(
			func(i int) { called = "present" },
			func() { called = "empty" },
		)
		if called != "present" {
			t.Error("expected 'present'")
		}

		None[int]().IfPresentOrElse(
			func(i int) { called = "present" },
			func() { called = "empty" },
		)
		if called != "empty" {
			t.Error("expected 'empty'")
		}
	})

	t.Run("OrElse", func(t *testing.T) {
		o := None[int]()
		res := o.OrElse(func() Optional[int] {
			return Some(42)
		})
		if !res.HasValue() || res.value != 42 {
			t.Errorf("expected Some(42), got %#v", res)
		}

		o2 := Some(100)
		res2 := o2.OrElse(func() Optional[int] {
			return Some(42)
		})
		if !res2.HasValue() || res2.value != 100 {
			t.Errorf("expected Some(100), got %#v", res2)
		}
	})
}
