package main

type Optional[T any] struct {
	value    T
	hasValue bool
}

func Some[T any](t T) Optional[T] {
	return Optional[T]{value: t, hasValue: true}
}

func None[T any]() Optional[T] {
	return Optional[T]{hasValue: false}
}

func (o Optional[T]) HasValue() bool {
	return o.hasValue
}

func (o Optional[T]) ValueOr(def T) T {
	if !o.hasValue {
		return def
	}
	return o.value
}

func (o Optional[T]) Get() (T, bool) {
	return o.value, o.hasValue
}

func (o *Optional[T]) Emplace(t T) *T {
	o.value = t
	o.hasValue = true
	return &o.value
}

func (o *Optional[T]) Reset() {
	o.value = *new(T)
	o.hasValue = false
}

func (o Optional[T]) AndThen[Z any](f func(T) Optional[Z]) Optional[Z] {
	if !o.hasValue {
		return Optional[Z]{hasValue: false}
	}
	return f(o.value)
}

func (o Optional[T]) Transform[Z any](f func(T) Z) Optional[Z] {
	if !o.hasValue {
		return Optional[Z]{hasValue: false}
	}
	return Optional[Z]{value: f(o.value), hasValue: true}
}

func (o Optional[T]) OrElse(f func() Optional[T]) Optional[T] {
	if o.hasValue {
		return o
	}
	return f()
}
