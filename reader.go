package main

// Reader represents a computation that depends on an environment of type R.
type Reader[R, T any] struct {
	run func(R) T
}

// NewReader returns a new Reader with the given computation.
func NewReader[R, T any](f func(R) T) Reader[R, T] {
	return Reader[R, T]{run: f}
}

// Run executes the Reader with the given environment.
func (r Reader[R, T]) Run(ctx R) T {
	return r.run(ctx)
}

// Map returns a new Reader that applies the mapping function to the result of this Reader.
func (r Reader[R, T]) Map[Z any](f func(T) Z) Reader[R, Z] {
	return NewReader(func(ctx R) Z {
		return f(r.Run(ctx))
	})
}

// AndThen returns a new Reader that applies the mapping function (returning another Reader) to the result of this Reader.
func (r Reader[R, T]) AndThen[Z any](f func(T) Reader[R, Z]) Reader[R, Z] {
	return NewReader(func(ctx R) Z {
		return f(r.Run(ctx)).Run(ctx)
	})
}

// Ask returns a Reader that simply returns the environment.
func Ask[R any]() Reader[R, R] {
	return NewReader(func(ctx R) R {
		return ctx
	})
}
