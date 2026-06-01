package main

import (
	"errors"
	"testing"
)

func TestResult(t *testing.T) {
	t.Run("Ok and Err", func(t *testing.T) {
		r := Ok(42)
		if !r.IsOk() || r.IsErr() || r.value != 42 {
			t.Error("expected Ok(42)")
		}

		err := errors.New("fail")
		re := Err[int](err)
		if re.IsOk() || !re.IsErr() || re.err != err {
			t.Error("expected Err")
		}
	})

	t.Run("From", func(t *testing.T) {
		r := From(42, nil)
		if !r.IsOk() || r.value != 42 {
			t.Error("expected Ok(42) from From")
		}

		err := errors.New("fail")
		re := From(0, err)
		if re.IsOk() || re.err != err {
			t.Error("expected Err from From")
		}
	})

	t.Run("ValueOr", func(t *testing.T) {
		if Ok(42).ValueOr(100) != 42 {
			t.Error("expected 42")
		}
		if Err[int](errors.New("fail")).ValueOr(100) != 100 {
			t.Error("expected 100")
		}
	})

	t.Run("Get", func(t *testing.T) {
		v, err := Ok(42).Get()
		if err != nil || v != 42 {
			t.Error("expected 42, nil")
		}

		v, err = Err[int](errors.New("fail")).Get()
		if err == nil || v != 0 {
			t.Error("expected 0, error")
		}
	})

	t.Run("AndThen", func(t *testing.T) {
		r := Ok(21).AndThen(func(i int) Result[int] {
			return Ok(i * 2)
		})
		if r.value != 42 {
			t.Error("expected 42")
		}

		re := Err[int](errors.New("fail")).AndThen(func(i int) Result[int] {
			return Ok(i * 2)
		})
		if re.IsOk() {
			t.Error("expected Err")
		}
	})

	t.Run("Map", func(t *testing.T) {
		r := Ok(42).Map(func(i int) string {
			return "success"
		})
		if r.value != "success" {
			t.Error("expected 'success'")
		}
	})

	t.Run("MapErr", func(t *testing.T) {
		r := Err[int](errors.New("fail")).MapErr(func(e error) error {
			return errors.New("mapped fail")
		})
		if r.err.Error() != "mapped fail" {
			t.Error("expected mapped error")
		}
	})

	t.Run("OrElse", func(t *testing.T) {
		r := Err[int](errors.New("fail")).OrElse(func(e error) Result[int] {
			return Ok(42)
		})
		if r.value != 42 {
			t.Error("expected 42")
		}
	})
}
