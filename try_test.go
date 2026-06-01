package main

import (
	"errors"
	"testing"
)

func TestTry(t *testing.T) {
	t.Run("Success and Failure", func(t *testing.T) {
		s := SuccessTry(42)
		if !s.IsSuccess() || s.IsFailure() || s.value != 42 {
			t.Error("expected SuccessTry")
		}

		err := errors.New("fail")
		f := FailureTry[int](err)
		if f.IsSuccess() || !f.IsFailure() || f.err != err {
			t.Error("expected FailureTry")
		}
	})

	t.Run("Invoke - Success", func(t *testing.T) {
		res := Invoke(func() (int, error) {
			return 42, nil
		})
		if !res.IsSuccess() || res.value != 42 {
			t.Error("expected SuccessTry(42)")
		}
	})

	t.Run("Invoke - Error", func(t *testing.T) {
		err := errors.New("fail")
		res := Invoke(func() (int, error) {
			return 0, err
		})
		if !res.IsFailure() || res.err != err {
			t.Error("expected FailureTry")
		}
	})

	t.Run("Invoke - Panic", func(t *testing.T) {
		res := Invoke(func() (int, error) {
			panic("something went wrong")
		})
		if !res.IsFailure() || res.err.Error() != "panic: something went wrong" {
			t.Error("expected FailureTry from panic")
		}
	})

	t.Run("Map", func(t *testing.T) {
		s := SuccessTry(42).Map(func(i int) string {
			return "success"
		})
		if s.value != "success" {
			t.Error("expected 'success'")
		}

		f := FailureTry[int](errors.New("fail")).Map(func(i int) string {
			return "success"
		})
		if f.IsSuccess() {
			t.Error("expected FailureTry")
		}
	})

	t.Run("Recover", func(t *testing.T) {
		res := FailureTry[int](errors.New("fail")).Recover(func(e error) int {
			return 42
		})
		if res.value != 42 {
			t.Error("expected 42")
		}
	})

	t.Run("RecoverWith", func(t *testing.T) {
		res := FailureTry[int](errors.New("fail")).RecoverWith(func(e error) Try[int] {
			return SuccessTry(42)
		})
		if res.value != 42 {
			t.Error("expected 42")
		}
	})
}
