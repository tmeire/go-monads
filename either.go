package main

type Either[L, R any] struct {
	left   L
	right  R
	isLeft bool
}

func Left[L, R any](l L) Either[L, R] {
	return Either[L, R]{left: l, isLeft: true}
}

func Right[L, R any](r R) Either[L, R] {
	return Either[L, R]{right: r, isLeft: false}
}

func (e Either[L, R]) IsLeft() bool {
	return e.isLeft
}

func (e Either[L, R]) IsRight() bool {
	return !e.isLeft
}

func (e Either[L, R]) GetLeft() (L, bool) {
	return e.left, e.isLeft
}

func (e Either[L, R]) GetRight() (R, bool) {
	return e.right, !e.isLeft
}

func (e Either[L, R]) MapLeft[NL any](f func(L) NL) Either[NL, R] {
	if e.isLeft {
		return Left[NL, R](f(e.left))
	}
	return Right[NL, R](e.right)
}

func (e Either[L, R]) MapRight[NR any](f func(R) NR) Either[L, NR] {
	if !e.isLeft {
		return Right[L, NR](f(e.right))
	}
	return Left[L, NR](e.left)
}

func (e Either[L, R]) Match[Z any](l func(L) Z, r func(R) Z) Z {
	if e.isLeft {
		return l(e.left)
	}
	return r(e.right)
}
