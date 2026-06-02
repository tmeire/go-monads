package main

import (
	"fmt"
	"testing"
)

type Config struct {
	BaseURL string
}

func TestReader(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		r := NewReader(func(cfg Config) string {
			return cfg.BaseURL + "/api"
		})
		cfg := Config{BaseURL: "https://example.com"}
		got := r.Run(cfg)
		want := "https://example.com/api"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Map", func(t *testing.T) {
		r := NewReader(func(cfg Config) string {
			return cfg.BaseURL
		}).Map(func(s string) int {
			return len(s)
		})
		cfg := Config{BaseURL: "abc"}
		got := r.Run(cfg)
		if got != 3 {
			t.Errorf("got %v, want 3", got)
		}
	})

	t.Run("AndThen", func(t *testing.T) {
		r := NewReader(func(cfg Config) string {
			return cfg.BaseURL
		}).AndThen(func(s string) Reader[Config, string] {
			return NewReader(func(cfg Config) string {
				return fmt.Sprintf("%s (%s)", s, cfg.BaseURL)
			})
		})
		cfg := Config{BaseURL: "host"}
		got := r.Run(cfg)
		want := "host (host)"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Ask", func(t *testing.T) {
		r := Ask[Config]().Map(func(cfg Config) string {
			return cfg.BaseURL
		})
		cfg := Config{BaseURL: "test"}
		got := r.Run(cfg)
		if got != "test" {
			t.Errorf("got %q, want test", got)
		}
	})
}
