package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestValidation(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		v := Valid[int, string](42)
		if !v.IsValid() {
			t.Error("Expected valid")
		}
		if v.Value() != 42 {
			t.Errorf("Expected 42, got %v", v.Value())
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		v := Invalid[int, string]("err1", "err2")
		if v.IsValid() {
			t.Error("Expected invalid")
		}
		expectedErrors := []string{"err1", "err2"}
		if !reflect.DeepEqual(v.Errors(), expectedErrors) {
			t.Errorf("Expected %v, got %v", expectedErrors, v.Errors())
		}
	})

	t.Run("Map", func(t *testing.T) {
		v := Valid[int, string](21).Map(func(i int) int { return i * 2 })
		if v.Value() != 42 {
			t.Errorf("Expected 42, got %v", v.Value())
		}

		v2 := Invalid[int, string]("err").Map(func(i int) int { return i * 2 })
		if v2.IsValid() {
			t.Error("Expected invalid")
		}
	})

	t.Run("And", func(t *testing.T) {
		v1 := Valid[int, string](10)
		v2 := Valid[int, string](20)
		combined := v1.And(v2, func(a, b int) int { return a + b })
		if combined.Value() != 30 {
			t.Errorf("Expected 30, got %v", combined.Value())
		}

		v3 := Invalid[int, string]("err1")
		v4 := Invalid[int, string]("err2")
		combined2 := v3.And(v4, func(a, b int) int { return a + b })
		if combined2.IsValid() {
			t.Error("Expected invalid")
		}
		expectedErrors := []string{"err1", "err2"}
		if !reflect.DeepEqual(combined2.Errors(), expectedErrors) {
			t.Errorf("Expected %v, got %v", expectedErrors, combined2.Errors())
		}

		combined3 := v1.And(v3, func(a, b int) int { return a + b })
		if combined3.IsValid() {
			t.Error("Expected invalid")
		}
		if !reflect.DeepEqual(combined3.Errors(), []string{"err1"}) {
			t.Errorf("Expected [err1], got %v", combined3.Errors())
		}
	})

	t.Run("ToResult", func(t *testing.T) {
		v := Valid[int, error](42)
		r := ToResult(v)
		if !r.IsOk() {
			t.Error("Expected Ok")
		}

		err := errors.New("fail")
		v2 := Invalid[int, error](err)
		r2 := ToResult(v2)
		if !r2.IsErr() {
			t.Error("Expected Err")
		}
		_, gotErr := r2.Get()
		if gotErr != err {
			t.Errorf("Expected %v, got %v", err, gotErr)
		}
	})
}
