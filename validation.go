package main

// Validation represents a value that is either valid (containing a value of type T)
// or invalid (containing one or more errors of type E).
// Unlike Result, Validation is designed to accumulate errors.
type Validation[T any, E any] struct {
	value   T
	errors  []E
	isValid bool
}

// Valid returns a Validation containing a success value.
func Valid[T, E any](t T) Validation[T, E] {
	return Validation[T, E]{value: t, isValid: true}
}

// Invalid returns a Validation containing one or more errors.
func Invalid[T, E any](errs ...E) Validation[T, E] {
	return Validation[T, E]{errors: errs, isValid: false}
}

// IsValid returns true if the Validation is valid.
func (v Validation[T, E]) IsValid() bool {
	return v.isValid
}

// Errors returns the accumulated errors.
func (v Validation[T, E]) Errors() []E {
	return v.errors
}

// Value returns the value if valid, otherwise the zero value of T.
func (v Validation[T, E]) Value() T {
	return v.value
}

// Map applies a function to the valid value.
func (v Validation[T, E]) Map[Z any](f func(T) Z) Validation[Z, E] {
	if !v.isValid {
		return Validation[Z, E]{errors: v.errors, isValid: false}
	}
	return Valid[Z, E](f(v.value))
}

// And combines two Validations. If both are valid, it returns a new Validation
// produced by the combiner function. If either or both are invalid, it accumulates all errors.
func (v Validation[T, E]) And[U, Z any](other Validation[U, E], f func(T, U) Z) Validation[Z, E] {
	if v.isValid && other.isValid {
		return Valid[Z, E](f(v.value, other.value))
	}
	var errs []E
	errs = append(errs, v.errors...)
	errs = append(errs, other.errors...)
	return Invalid[Z, E](errs...)
}

// ToResult converts Validation to a Result if the error type is error.
// Note: This only takes the first error if invalid.
func ToResult[T any](v Validation[T, error]) Result[T] {
	if !v.isValid {
		if len(v.errors) > 0 {
			return Err[T](v.errors[0])
		}
		return Err[T](nil)
	}
	return Ok(v.value)
}
