package main

import (
	"sync/atomic"
	"testing"
)

func TestLazy(t *testing.T) {
	t.Run("Initialization", func(t *testing.T) {
		var count atomic.Int32
		l := NewLazy(func() int {
			count.Add(1)
			return 42
		})

		if count.Load() != 0 {
			t.Error("expected 0 calls before Get")
		}

		if val := l.Get(); val != 42 {
			t.Error("expected 42")
		}

		if val := l.Get(); val != 42 {
			t.Error("expected 42 on second call")
		}

		if count.Load() != 1 {
			t.Error("expected exactly 1 call to init")
		}
	})

	t.Run("Map", func(t *testing.T) {
		var count atomic.Int32
		l := NewLazy(func() int {
			count.Add(1)
			return 21
		}).Map(func(i int) string {
			return "success"
		})

		if count.Load() != 0 {
			t.Error("expected 0 calls before Get")
		}

		if val := l.Get(); val != "success" {
			t.Error("expected 'success'")
		}
	})

	t.Run("AndThen", func(t *testing.T) {
		l := NewLazy(func() int {
			return 21
		}).AndThen(func(i int) *Lazy[int] {
			return NewLazy(func() int {
				return i * 2
			})
		})

		if val := l.Get(); val != 42 {
			t.Error("expected 42")
		}
	})
}
