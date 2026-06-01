package main

import "iter"

type Stream[T any] struct {
	seq iter.Seq[T]
}

func NewStream[T any](seq iter.Seq[T]) Stream[T] {
	return Stream[T]{seq: seq}
}

func OfSlice[T any](slice []T) Stream[T] {
	return Stream[T]{seq: func(yield func(T) bool) {
		for _, v := range slice {
			if !yield(v) {
				return
			}
		}
	}}
}

func (s Stream[T]) Filter(f func(T) bool) Stream[T] {
	return Stream[T]{seq: func(yield func(T) bool) {
		for v := range s.seq {
			if f(v) {
				if !yield(v) {
					return
				}
			}
		}
	}}
}

func (s Stream[T]) Map[Z any](f func(T) Z) Stream[Z] {
	return Stream[Z]{seq: func(yield func(Z) bool) {
		for v := range s.seq {
			if !yield(f(v)) {
				return
			}
		}
	}}
}

func (s Stream[T]) Reduce(init T, f func(T, T) T) T {
	acc := init
	for v := range s.seq {
		acc = f(acc, v)
	}
	return acc
}

func (s Stream[T]) ToSlice() []T {
	var res []T
	for v := range s.seq {
		res = append(res, v)
	}
	return res
}

func (s Stream[T]) FindFirst() Optional[T] {
	for v := range s.seq {
		return Some(v)
	}
	return None[T]()
}
