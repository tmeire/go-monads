package main

// Optional represents a container object which may or may not contain a non-nil value.
type Optional[T any] struct {
	value    T
	hasValue bool
}

// Some returns an Optional describing the given non-nil value.
func Some[T any](t T) Optional[T] {
	return Optional[T]{value: t, hasValue: true}
}

// None returns an empty Optional instance. No value is present for this Optional.
func None[T any]() Optional[T] {
	return Optional[T]{hasValue: false}
}

// HasValue returns true if there is a value present, otherwise false.
func (o Optional[T]) HasValue() bool {
	return o.hasValue
}

// ValueOr returns the value if present, otherwise returns def.
func (o Optional[T]) ValueOr(def T) T {
	if !o.hasValue {
		return def
	}
	return o.value
}

// ValueOrGet returns the value if present, otherwise invokes the producer function f and returns its result.
func (o Optional[T]) ValueOrGet(f func() T) T {
	if !o.hasValue {
		return f()
	}
	return o.value
}

// ValueOrErr returns the value and nil if present, otherwise returns a zero value and the provided error.
func (o Optional[T]) ValueOrErr(err error) (T, error) {
	if !o.hasValue {
		return *new(T), err
	}
	return o.value, nil
}

// Get returns the value and a boolean indicating presence.
func (o Optional[T]) Get() (T, bool) {
	return o.value, o.hasValue
}

// Filter if a value is present, and the value matches the given predicate, returns an Optional describing the value, otherwise returns an empty Optional.
func (o Optional[T]) Filter(f func(T) bool) Optional[T] {
	if !o.hasValue || !f(o.value) {
		return Optional[T]{hasValue: false}
	}
	return o
}

// IfPresent if a value is present, performs the given action with the value, otherwise does nothing.
func (o Optional[T]) IfPresent(f func(T)) {
	if o.hasValue {
		f(o.value)
	}
}

// IfPresentOrElse if a value is present, performs the given action with the value, otherwise performs the given empty-based action.
func (o Optional[T]) IfPresentOrElse(f func(T), e func()) {
	if o.hasValue {
		f(o.value)
	} else {
		e()
	}
}

// Emplace constructs the value in-place, resetting the Optional to a success state.
// It returns a pointer to the internal value for direct modification.
func (o *Optional[T]) Emplace(t T) *T {
	o.value = t
	o.hasValue = true
	return &o.value
}

// Reset clears the Optional, making it empty.
func (o *Optional[T]) Reset() {
	o.value = *new(T)
	o.hasValue = false
}

// AndThen if a value is present, returns the result of applying the given Optional-bearing mapping function to the value, otherwise returns an empty Optional.
func (o Optional[T]) AndThen[Z any](f func(T) Optional[Z]) Optional[Z] {
	if !o.hasValue {
		return Optional[Z]{hasValue: false}
	}
	return f(o.value)
}

// Transform if a value is present, returns an Optional describing the result of applying the given mapping function to the value, otherwise returns an empty Optional.
func (o Optional[T]) Transform[Z any](f func(T) Z) Optional[Z] {
	if !o.hasValue {
		return Optional[Z]{hasValue: false}
	}
	return Optional[Z]{value: f(o.value), hasValue: true}
}

// OrElse if a value is present, returns the Optional, otherwise returns an Optional produced by the supplying function.
func (o Optional[T]) OrElse(f func() Optional[T]) Optional[T] {
	if o.hasValue {
		return o
	}
	return f()
}
