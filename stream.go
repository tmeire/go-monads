package main

import "iter"

// Stream represents a lazy sequence of elements.
// It is a wrapper around iter.Seq[T].
type Stream[T any] struct {
	seq iter.Seq[T]
}

// NewStream creates a new Stream from an existing iterator.
func NewStream[T any](seq iter.Seq[T]) Stream[T] {
	return Stream[T]{seq: seq}
}

// OfSlice creates a new Stream from the given slice.
func OfSlice[T any](slice []T) Stream[T] {
	return Stream[T]{seq: func(yield func(T) bool) {
		for _, v := range slice {
			if !yield(v) {
				return
			}
		}
	}}
}

// Filter returns a new Stream consisting of the elements of this stream that match the given predicate.
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

// Map returns a new Stream consisting of the results of applying the given function to the elements of this stream.
func (s Stream[T]) Map[Z any](f func(T) Z) Stream[Z] {
	return Stream[Z]{seq: func(yield func(Z) bool) {
		for v := range s.seq {
			if !yield(f(v)) {
				return
			}
		}
	}}
}

// Reduce performs a reduction on the elements of this stream, using the provided initial value and an associative accumulation function.
func (s Stream[T]) Reduce(init T, f func(T, T) T) T {
	acc := init
	for v := range s.seq {
		acc = f(acc, v)
	}
	return acc
}

// ToSlice collects the elements of the stream into a slice.
func (s Stream[T]) ToSlice() []T {
	var res []T
	for v := range s.seq {
		res = append(res, v)
	}
	return res
}

// FindFirst returns an Optional describing the first element of this stream, or an empty Optional if the stream is empty.
func (s Stream[T]) FindFirst() Optional[T] {
	for v := range s.seq {
		return Some(v)
	}
	return None[T]()
}
