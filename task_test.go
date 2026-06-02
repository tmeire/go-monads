package main

import (
	"errors"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		task := NewTask(func() Result[int] {
			return Ok(42)
		})
		res := task.Run()
		if !res.IsOk() {
			t.Error("Expected Ok")
		}
		val, _ := res.Get()
		if val != 42 {
			t.Errorf("got %v, want 42", val)
		}
	})

	t.Run("Map", func(t *testing.T) {
		task := NewTask(func() Result[int] {
			return Ok(21)
		}).Map(func(i int) int { return i * 2 })
		res := task.Run()
		val, _ := res.Get()
		if val != 42 {
			t.Errorf("got %v, want 42", val)
		}
	})

	t.Run("AndThen", func(t *testing.T) {
		task := NewTask(func() Result[int] {
			return Ok(10)
		}).AndThen(func(i int) Task[int] {
			return NewTask(func() Result[int] {
				return Ok(i + 32)
			})
		})
		res := task.Run()
		val, _ := res.Get()
		if val != 42 {
			t.Errorf("got %v, want 42", val)
		}

		// Error case
		err := errors.New("fail")
		taskErr := NewTask(func() Result[int] {
			return Err[int](err)
		}).AndThen(func(i int) Task[int] {
			return NewTask(func() Result[int] {
				return Ok(i + 32)
			})
		})
		resErr := taskErr.Run()
		if !resErr.IsErr() {
			t.Error("Expected Err")
		}
	})

	t.Run("RunAsync", func(t *testing.T) {
		start := time.Now()
		task := NewTask(func() Result[int] {
			time.Sleep(100 * time.Millisecond)
			return Ok(42)
		})
		ch := task.RunAsync()
		
		// Should be non-blocking
		if time.Since(start) > 50*time.Millisecond {
			t.Error("RunAsync blocked")
		}

		res := <-ch
		val, _ := res.Get()
		if val != 42 {
			t.Errorf("got %v, want 42", val)
		}
	})
}
