package main

type Expected[T any, E any] struct {
	t        T
	e        E
	hasValue bool
}

func Success[T, E any](t T) Expected[T, E] {
	return Expected[T, E]{t: t, hasValue: true}
}

func Failure[T, E any](e E) Expected[T, E] {
	return Expected[T, E]{e: e, hasValue: false}
}

func (e Expected[T, E]) HasValue() bool {
	return e.hasValue
}

func (e Expected[T, E]) ValueOr(def T) T {
	if !e.hasValue {
		return def
	}
	return e.t
}

func (e Expected[T, E]) ErrorOr(def E) E {
	if e.hasValue {
		return def
	}
	return e.e
}

func (e Expected[T, E]) Get() (T, E) {
	return e.t, e.e
}

func (e Expected[T, E]) AndThen[Z any](f func(T) Expected[Z, E]) Expected[Z, E] {
	if !e.hasValue {
		return Expected[Z, E]{e: e.e, hasValue: false}
	}
	return f(e.t)
}

func (e Expected[T, E]) Transform[Z any](f func(T) Z) Expected[Z, E] {
	if !e.hasValue {
		return Expected[Z, E]{e: e.e, hasValue: false}
	}
	return Expected[Z, E]{t: f(e.t), hasValue: true}
}

func (e Expected[T, E]) OrElse[Z any](f func(E) Expected[T, Z]) Expected[T, Z] {
	if e.hasValue {
		return Expected[T, Z]{t: e.t, hasValue: true}
	}
	return f(e.e)
}

func (e Expected[T, E]) TransformError[Z any](f func(E) Z) Expected[T, Z] {
	if e.hasValue {
		return Expected[T, Z]{t: e.t, hasValue: true}
	}
	return Expected[T, Z]{e: f(e.e), hasValue: false}
}

func (e *Expected[T, E]) Emplace(t T) *T {
	e.t = t
	e.e = *new(E)
	e.hasValue = true
	return &e.t
}

func (e *Expected[T, E]) EmplaceError(err E) *E {
	e.t = *new(T)
	e.e = err
	e.hasValue = false
	return &e.e
}
