package main

type Result[T any] struct {
	value T
	err   error
}

func Ok[T any](t T) Result[T] {
	return Result[T]{value: t, err: nil}
}

func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

func From[T any](t T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(t)
}

func (r Result[T]) IsOk() bool {
	return r.err == nil
}

func (r Result[T]) IsErr() bool {
	return r.err != nil
}

func (r Result[T]) ValueOr(def T) T {
	if r.err != nil {
		return def
	}
	return r.value
}

func (r Result[T]) ValueOrGet(f func() T) T {
	if r.err != nil {
		return f()
	}
	return r.value
}

func (r Result[T]) Get() (T, error) {
	return r.value, r.err
}

func (r Result[T]) AndThen[Z any](f func(T) Result[Z]) Result[Z] {
	if r.err != nil {
		return Err[Z](r.err)
	}
	return f(r.value)
}

func (r Result[T]) Map[Z any](f func(T) Z) Result[Z] {
	if r.err != nil {
		return Err[Z](r.err)
	}
	return Ok(f(r.value))
}

func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if r.err != nil {
		return Err[T](f(r.err))
	}
	return r
}

func (r Result[T]) OrElse(f func(error) Result[T]) Result[T] {
	if r.err != nil {
		return f(r.err)
	}
	return r
}
