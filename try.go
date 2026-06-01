package main

import (
	"fmt"
)

// Try represents a computation that may either result in an exception (panic), an error, or a successfully computed value.
type Try[T any] struct {
	value T
	err   error
}

// SuccessTry returns a Try describing a successful computation.
func SuccessTry[T any](t T) Try[T] {
	return Try[T]{value: t, err: nil}
}

// FailureTry returns a Try describing a failed computation.
func FailureTry[T any](err error) Try[T] {
	return Try[T]{err: err}
}

// Invoke executes the given function and returns a Try describing the result.
// It automatically recovers from panics and converts them into a FailureTry.
func Invoke[T any](f func() (T, error)) (res Try[T]) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				res = FailureTry[T](err)
			} else {
				res = FailureTry[T](fmt.Errorf("panic: %v", r))
			}
		}
	}()

	val, err := f()
	if err != nil {
		return FailureTry[T](err)
	}
	return SuccessTry(val)
}

// IsSuccess returns true if the Try is a success.
func (t Try[T]) IsSuccess() bool {
	return t.err == nil
}

// IsFailure returns true if the Try is a failure.
func (t Try[T]) IsFailure() bool {
	return t.err != nil
}

// Get returns both the value and the error.
func (t Try[T]) Get() (T, error) {
	return t.value, t.err
}

// Map if the Try is a success, returns a Try describing the result of applying the given mapping function to the value, otherwise returns the original failure.
func (t Try[T]) Map[Z any](f func(T) Z) Try[Z] {
	if t.err != nil {
		return FailureTry[Z](t.err)
	}
	return Invoke(func() (Z, error) {
		return f(t.value), nil
	})
}

// AndThen if the Try is a success, returns the result of applying the given Try-bearing mapping function to the value, otherwise returns the original failure.
func (t Try[T]) AndThen[Z any](f func(T) Try[Z]) Try[Z] {
	if t.err != nil {
		return FailureTry[Z](t.err)
	}
	return f(t.value)
}

// Recover if the Try is a failure, returns a success Try produced by the given recovery function, otherwise returns the original success Try.
func (t Try[T]) Recover(f func(error) T) Try[T] {
	if t.err != nil {
		return SuccessTry(f(t.err))
	}
	return t
}

// RecoverWith if the Try is a failure, returns a Try produced by the given recovery function, otherwise returns the original success Try.
func (t Try[T]) RecoverWith(f func(error) Try[T]) Try[T] {
	if t.err != nil {
		return f(t.err)
	}
	return t
}
