package main

// Either represents a value of one of two possible types (a disjoint union).
// An instance of Either is either an instance of Left or Right.
type Either[L, R any] struct {
	left   L
	right  R
	isLeft bool
}

// Left constructs a Left Either instance with the given value.
func Left[L, R any](l L) Either[L, R] {
	return Either[L, R]{left: l, isLeft: true}
}

// Right constructs a Right Either instance with the given value.
func Right[L, R any](r R) Either[L, R] {
	return Either[L, R]{right: r, isLeft: false}
}

// IsLeft returns true if this is a Left instance.
func (e Either[L, R]) IsLeft() bool {
	return e.isLeft
}

// IsRight returns true if this is a Right instance.
func (e Either[L, R]) IsRight() bool {
	return !e.isLeft
}

// GetLeft returns the left value and true if this is a Left instance, otherwise returns zero value and false.
func (e Either[L, R]) GetLeft() (L, bool) {
	return e.left, e.isLeft
}

// GetRight returns the right value and true if this is a Right instance, otherwise returns zero value and false.
func (e Either[L, R]) GetRight() (R, bool) {
	return e.right, !e.isLeft
}

// MapLeft returns a new Either after applying the mapping function if this is a Left instance.
func (e Either[L, R]) MapLeft[NL any](f func(L) NL) Either[NL, R] {
	if e.isLeft {
		return Left[NL, R](f(e.left))
	}
	return Right[NL, R](e.right)
}

// MapRight returns a new Either after applying the mapping function if this is a Right instance.
func (e Either[L, R]) MapRight[NR any](f func(R) NR) Either[L, NR] {
	if !e.isLeft {
		return Right[L, NR](f(e.right))
	}
	return Left[L, NR](e.left)
}

// Match applies the appropriate function based on the instance type and returns the result.
func (e Either[L, R]) Match[Z any](l func(L) Z, r func(R) Z) Z {
	if e.isLeft {
		return l(e.left)
	}
	return r(e.right)
}
