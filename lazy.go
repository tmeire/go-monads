package main

import "sync"

// Lazy represents a value that is only computed when first accessed.
// It is thread-safe and caches the result for subsequent accesses.
type Lazy[T any] struct {
	once  sync.Once
	value T
	init  func() T
}

// NewLazy returns a new Lazy instance with the given initialization function.
func NewLazy[T any](init func() T) *Lazy[T] {
	return &Lazy[T]{init: init}
}

// Get returns the computed value, initializing it if necessary.
func (l *Lazy[T]) Get() T {
	l.once.Do(func() {
		l.value = l.init()
	})
	return l.value
}

// Map returns a new Lazy instance that will apply the given mapping function to the result of this Lazy.
func (l *Lazy[T]) Map[Z any](f func(T) Z) *Lazy[Z] {
	return NewLazy(func() Z {
		return f(l.Get())
	})
}

// AndThen returns a new Lazy instance that will apply the given mapping function (returning another Lazy) to the result of this Lazy.
func (l *Lazy[T]) AndThen[Z any](f func(T) *Lazy[Z]) *Lazy[Z] {
	return NewLazy(func() Z {
		return f(l.Get()).Get()
	})
}
