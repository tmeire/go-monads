package main

type Expected[T comparable, E comparable] struct {
	t T
	e E
}

func (e Expected[T, E]) AndThen[Z comparable](f func(T) Expected[Z, E]) Expected[Z, E] {
	var z E
	if e.e != z {
		return Expected[Z, E]{e: e.e}
	}
	return f(e.t)
}

func (e Expected[T, E]) Transform[Z comparable](f func(T) Z) Expected[Z, E] {
	var z E
	if e.e != z {
		return Expected[Z, E]{e: e.e}
	}
	return Expected[Z, E]{t: f(e.t)}
}

func (e Expected[T, E]) OrElse[Z comparable](f func(E) Expected[T, Z]) Expected[T, Z] {
	var z T
	if e.t != z {
		return Expected[T, Z]{t: e.t}
	}
	return f(e.e)
}

func (e Expected[T, E]) TransformError[Z comparable](f func(E) Z) Expected[T, Z] {
	var z T
	if e.t != z {
		return Expected[T, Z]{t: e.t}
	}
	return Expected[T, Z]{e: f(e.e)}
}
