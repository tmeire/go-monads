package main

import (
	"reflect"
	"testing"
)

func TestWriter(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		w := NewWriter(42, "init")
		if w.Value() != 42 {
			t.Errorf("got %v, want 42", w.Value())
		}
		if !reflect.DeepEqual(w.Logs(), []string{"init"}) {
			t.Errorf("got %v, want [init]", w.Logs())
		}
	})

	t.Run("Map", func(t *testing.T) {
		w := NewWriter(21, "start").Map(func(i int) int { return i * 2 })
		if w.Value() != 42 {
			t.Errorf("got %v, want 42", w.Value())
		}
		if !reflect.DeepEqual(w.Logs(), []string{"start"}) {
			t.Errorf("got %v, want [start]", w.Logs())
		}
	})

	t.Run("AndThen", func(t *testing.T) {
		w := NewWriter(10, "step 1").AndThen(func(i int) Writer[string, int] {
			return NewWriter(i+20, "step 2")
		})
		if w.Value() != 30 {
			t.Errorf("got %v, want 30", w.Value())
		}
		expectedLogs := []string{"step 1", "step 2"}
		if !reflect.DeepEqual(w.Logs(), expectedLogs) {
			t.Errorf("got %v, want %v", w.Logs(), expectedLogs)
		}
	})

	t.Run("Tell", func(t *testing.T) {
		w := Tell("log message")
		if w.Value() != nil {
			t.Errorf("got %v, want nil", w.Value())
		}
		if !reflect.DeepEqual(w.Logs(), []string{"log message"}) {
			t.Errorf("got %v, want [log message]", w.Logs())
		}
	})
}
