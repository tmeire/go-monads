package main

import "sync"

type Lazy[T any] struct {
	once  sync.Once
	value T
	init  func() T
}

func NewLazy[T any](init func() T) *Lazy[T] {
	return &Lazy[T]{init: init}
}

func (l *Lazy[T]) Get() T {
	l.once.Do(func() {
		l.value = l.init()
	})
	return l.value
}

func (l *Lazy[T]) Map[Z any](f func(T) Z) *Lazy[Z] {
	return NewLazy(func() Z {
		return f(l.Get())
	})
}

func (l *Lazy[T]) AndThen[Z any](f func(T) *Lazy[Z]) *Lazy[Z] {
	return NewLazy(func() Z {
		return f(l.Get()).Get()
	})
}
