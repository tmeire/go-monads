package main

import (
	"testing"
)

func TestEither(t *testing.T) {
	t.Run("Left and Right", func(t *testing.T) {
		l := Left[int, string](42)
		if !l.IsLeft() || l.IsRight() {
			t.Error("expected Left")
		}

		r := Right[int, string]("success")
		if !r.IsRight() || r.IsLeft() {
			t.Error("expected Right")
		}
	})

	t.Run("GetLeft and GetRight", func(t *testing.T) {
		l := Left[int, string](42)
		val, ok := l.GetLeft()
		if !ok || val != 42 {
			t.Error("expected 42, true")
		}
		str, ok := l.GetRight()
		if ok || str != "" {
			t.Error("expected empty string, false")
		}
	})

	t.Run("MapLeft", func(t *testing.T) {
		l := Left[int, string](42).MapLeft(func(i int) float64 {
			return float64(i) * 1.5
		})
		if val, _ := l.GetLeft(); val != 63.0 {
			t.Error("expected 63.0")
		}

		r := Right[int, string]("success").MapLeft(func(i int) float64 {
			return float64(i) * 1.5
		})
		if str, _ := r.GetRight(); str != "success" {
			t.Error("expected 'success'")
		}
	})

	t.Run("MapRight", func(t *testing.T) {
		r := Right[int, string]("success").MapRight(func(s string) int {
			return len(s)
		})
		if val, _ := r.GetRight(); val != 7 {
			t.Error("expected 7")
		}

		l := Left[int, string](42).MapRight(func(s string) int {
			return len(s)
		})
		if val, _ := l.GetLeft(); val != 42 {
			t.Error("expected 42")
		}
	})

	t.Run("Match", func(t *testing.T) {
		l := Left[int, string](42)
		resL := l.Match(
			func(i int) string { return "left" },
			func(s string) string { return "right" },
		)
		if resL != "left" {
			t.Error("expected 'left'")
		}

		r := Right[int, string]("success")
		resR := r.Match(
			func(i int) string { return "left" },
			func(s string) string { return "right" },
		)
		if resR != "right" {
			t.Error("expected 'right'")
		}
	})
}
