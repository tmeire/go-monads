package main

import (
	"fmt"
)

type Try[T any] struct {
	value T
	err   error
}

func SuccessTry[T any](t T) Try[T] {
	return Try[T]{value: t, err: nil}
}

func FailureTry[T any](err error) Try[T] {
	return Try[T]{err: err}
}

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

func (t Try[T]) IsSuccess() bool {
	return t.err == nil
}

func (t Try[T]) IsFailure() bool {
	return t.err != nil
}

func (t Try[T]) Get() (T, error) {
	return t.value, t.err
}

func (t Try[T]) Map[Z any](f func(T) Z) Try[Z] {
	if t.err != nil {
		return FailureTry[Z](t.err)
	}
	return Invoke(func() (Z, error) {
		return f(t.value), nil
	})
}

func (t Try[T]) AndThen[Z any](f func(T) Try[Z]) Try[Z] {
	if t.err != nil {
		return FailureTry[Z](t.err)
	}
	return f(t.value)
}

func (t Try[T]) Recover(f func(error) T) Try[T] {
	if t.err != nil {
		return SuccessTry(f(t.err))
	}
	return t
}

func (t Try[T]) RecoverWith(f func(error) Try[T]) Try[T] {
	if t.err != nil {
		return f(t.err)
	}
	return t
}
