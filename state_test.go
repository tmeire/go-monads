package main

import (
	"testing"
)

func TestState(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		s := NewState(func(i int) (string, int) {
			return "val", i + 1
		})
		val, next := s.Run(10)
		if val != "val" || next != 11 {
			t.Errorf("got (%q, %v), want (\"val\", 11)", val, next)
		}
	})

	t.Run("Map", func(t *testing.T) {
		s := NewState(func(i int) (int, int) {
			return i, i + 1
		}).Map(func(i int) int { return i * 2 })
		val, next := s.Run(10)
		if val != 20 || next != 11 {
			t.Errorf("got (%v, %v), want (20, 11)", val, next)
		}
	})

	t.Run("AndThen", func(t *testing.T) {
		s := NewState(func(i int) (int, int) {
			return i, i + 1
		}).AndThen(func(i int) State[int, int] {
			return NewState(func(state int) (int, int) {
				return i + state, state + 5
			})
		})
		// Run(10) -> (10, 11)
		// next: NewState(func(11) (10+11, 11+5)) -> (21, 16)
		val, next := s.Run(10)
		if val != 21 || next != 16 {
			t.Errorf("got (%v, %v), want (21, 16)", val, next)
		}
	})

	t.Run("Utilities", func(t *testing.T) {
		program := Get[int]().AndThen(func(s int) State[int, any] {
			return Put(s + 10)
		}).AndThen(func(_ any) State[int, int] {
			return Modify(func(s int) int { return s * 2 }).AndThen(func(_ any) State[int, int] {
				return Get[int]()
			})
		})

		// Get(5) -> (5, 5)
		// Put(15) -> (nil, 15)
		// Modify(*2) -> (nil, 30)
		// Get() -> (30, 30)
		val, next := program.Run(5)
		if val != 30 || next != 30 {
			t.Errorf("got (%v, %v), want (30, 30)", val, next)
		}
	})
}
