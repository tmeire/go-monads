package main

// Result represents a value that is either a success (Ok) or a failure (Err).
// It is a specialized version of Expected[T, error].
type Result[T any] struct {
	value T
	err   error
}

// Ok returns a Result describing a success with the given value.
func Ok[T any](t T) Result[T] {
	return Result[T]{value: t, err: nil}
}

// Err returns a Result describing a failure with the given error.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// From constructs a Result from a value and an error, following the standard Go return pattern.
func From[T any](t T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(t)
}

// IsOk returns true if the Result is a success.
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr returns true if the Result is a failure.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// ValueOr returns the success value if present, otherwise returns def.
func (r Result[T]) ValueOr(def T) T {
	if r.err != nil {
		return def
	}
	return r.value
}

// ValueOrGet returns the success value if present, otherwise invokes the producer function f and returns its result.
func (r Result[T]) ValueOrGet(f func() T) T {
	if r.err != nil {
		return f()
	}
	return r.value
}

// Get returns both the value and the error.
func (r Result[T]) Get() (T, error) {
	return r.value, r.err
}

// AndThen if the Result is a success, returns the result of applying the given Result-bearing mapping function to the value, otherwise returns the original error.
func (r Result[T]) AndThen[Z any](f func(T) Result[Z]) Result[Z] {
	if r.err != nil {
		return Err[Z](r.err)
	}
	return f(r.value)
}

// Map if the Result is a success, returns a Result describing the result of applying the given mapping function to the value, otherwise returns the original error.
func (r Result[T]) Map[Z any](f func(T) Z) Result[Z] {
	if r.err != nil {
		return Err[Z](r.err)
	}
	return Ok(f(r.value))
}

// MapErr if the Result is a failure, returns a Result describing the result of applying the given mapping function to the error, otherwise returns the original success value.
func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if r.err != nil {
		return Err[T](f(r.err))
	}
	return r
}

// OrElse if the Result is a failure, returns a Result produced by the given error-handling function, otherwise returns the original success value.
func (r Result[T]) OrElse(f func(error) Result[T]) Result[T] {
	if r.err != nil {
		return f(r.err)
	}
	return r
}
