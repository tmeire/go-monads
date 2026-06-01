package main

import (
	"reflect"
	"testing"
)

func TestStream(t *testing.T) {
	t.Run("OfSlice and ToSlice", func(t *testing.T) {
		input := []int{1, 2, 3}
		res := OfSlice(input).ToSlice()
		if !reflect.DeepEqual(res, input) {
			t.Errorf("expected %v, got %v", input, res)
		}
	})

	t.Run("Filter", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		res := OfSlice(input).Filter(func(i int) bool {
			return i%2 == 0
		}).ToSlice()
		expected := []int{2, 4}
		if !reflect.DeepEqual(res, expected) {
			t.Errorf("expected %v, got %v", expected, res)
		}
	})

	t.Run("Map", func(t *testing.T) {
		input := []int{1, 2}
		res := OfSlice(input).Map(func(i int) string {
			if i == 1 {
				return "one"
			}
			return "two"
		}).ToSlice()
		expected := []string{"one", "two"}
		if !reflect.DeepEqual(res, expected) {
			t.Errorf("expected %v, got %v", expected, res)
		}
	})

	t.Run("Reduce", func(t *testing.T) {
		input := []int{1, 2, 3}
		res := OfSlice(input).Reduce(10, func(acc int, val int) int {
			return acc + val
		})
		if res != 16 {
			t.Errorf("expected 16, got %d", res)
		}
	})

	t.Run("FindFirst", func(t *testing.T) {
		input := []int{1, 2, 3}
		opt := OfSlice(input).Filter(func(i int) bool { return i > 1 }).FindFirst()
		if !opt.HasValue() || opt.ValueOr(0) != 2 {
			t.Errorf("expected Some(2), got %v", opt)
		}

		optNone := OfSlice(input).Filter(func(i int) bool { return i > 5 }).FindFirst()
		if optNone.HasValue() {
			t.Errorf("expected None")
		}
	})
}
