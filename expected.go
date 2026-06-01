package main

import "reflect"

// Expected represents either of two values: an expected value of type T, or an unexpected value of type E.
type Expected[T any, E any] struct {
	t        T
	e        E
	hasValue bool
}

// Success returns an Expected containing the given success value.
func Success[T, E any](t T) Expected[T, E] {
	return Expected[T, E]{t: t, hasValue: true}
}

// Failure returns an Expected containing the given error value.
func Failure[T, E any](e E) Expected[T, E] {
	return Expected[T, E]{e: e, hasValue: false}
}

// HasValue returns true if the object contains a success value, otherwise false.
func (e Expected[T, E]) HasValue() bool {
	return e.hasValue
}

// ValueOr returns the value if present, otherwise returns def.
func (e Expected[T, E]) ValueOr(def T) T {
	if !e.hasValue {
		return def
	}
	return e.t
}

// ErrorOr returns the error if present, otherwise returns def.
func (e Expected[T, E]) ErrorOr(def E) E {
	if e.hasValue {
		return def
	}
	return e.e
}

// Get returns both the value and error. One will be its zero value.
func (e Expected[T, E]) Get() (T, E) {
	return e.t, e.e
}

// AndThen if a success value is present, returns the result of applying the given Expected-bearing mapping function to the value, otherwise returns the original error.
func (e Expected[T, E]) AndThen[Z any](f func(T) Expected[Z, E]) Expected[Z, E] {
	if !e.hasValue {
		return Expected[Z, E]{e: e.e, hasValue: false}
	}
	return f(e.t)
}

// Transform if a success value is present, returns an Expected describing the result of applying the given mapping function to the value, otherwise returns the original error.
func (e Expected[T, E]) Transform[Z any](f func(T) Z) Expected[Z, E] {
	if !e.hasValue {
		return Expected[Z, E]{e: e.e, hasValue: false}
	}
	return Expected[Z, E]{t: f(e.t), hasValue: true}
}

// OrElse if a success value is present, returns the Expected, otherwise returns an Expected produced by the given error-handling function.
func (e Expected[T, E]) OrElse[Z any](f func(E) Expected[T, Z]) Expected[T, Z] {
	if e.hasValue {
		return Expected[T, Z]{t: e.t, hasValue: true}
	}
	return f(e.e)
}

// TransformError if an error is present, returns an Expected describing the result of applying the given mapping function to the error, otherwise returns the original success value.
func (e Expected[T, E]) TransformError[Z any](f func(E) Z) Expected[T, Z] {
	if e.hasValue {
		return Expected[T, Z]{t: e.t, hasValue: true}
	}
	return Expected[T, Z]{e: f(e.e), hasValue: false}
}

// Emplace constructs the success value in-place, resetting any previous state.
// It returns a pointer to the internal value for direct modification.
func (e *Expected[T, E]) Emplace(t T) *T {
	e.t = t
	e.e = *new(E)
	e.hasValue = true
	return &e.t
}

// EmplaceError constructs the error value in-place, resetting any previous state.
// It returns a pointer to the internal error for direct modification.
func (e *Expected[T, E]) EmplaceError(err E) *E {
	e.t = *new(T)
	e.e = err
	e.hasValue = false
	return &e.e
}

// ValueOrGet returns the success value if present, otherwise invokes the producer function f and returns its result.
func (e Expected[T, E]) ValueOrGet(f func() T) T {
	if !e.hasValue {
		return f()
	}
	return e.t
}

// IfHasValue if a success value is present, performs the given action with the value, otherwise does nothing.
func (e Expected[T, E]) IfHasValue(f func(T)) {
	if e.hasValue {
		f(e.t)
	}
}

// IfHasError if an error is present, performs the given action with the error, otherwise does nothing.
func (e Expected[T, E]) IfHasError(f func(E)) {
	if !e.hasValue {
		f(e.e)
	}
}

// IfHasValueOrElse if a success value is present, performs the success action, otherwise performs the error action.
func (e Expected[T, E]) IfHasValueOrElse(f func(T), g func(E)) {
	if e.hasValue {
		f(e.t)
	} else {
		g(e.e)
	}
}

// Equals returns true if the other Expected has the same state and values (using DeepEqual).
func (e Expected[T, E]) Equals(other Expected[T, E]) bool {
	if e.hasValue != other.hasValue {
		return false
	}
	if e.hasValue {
		return reflect.DeepEqual(e.t, other.t)
	}
	return reflect.DeepEqual(e.e, other.e)
}
