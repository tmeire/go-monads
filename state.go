package main

// State represents a computation that carries a state of type S and produces a result of type T.
type State[S, T any] struct {
	run func(S) (T, S)
}

// NewState returns a new State with the given computation.
func NewState[S, T any](f func(S) (T, S)) State[S, T] {
	return State[S, T]{run: f}
}

// Run executes the State computation with the given initial state.
func (s State[S, T]) Run(initial S) (T, S) {
	return s.run(initial)
}

// Map returns a new State that applies the mapping function to the result.
func (s State[S, T]) Map[Z any](f func(T) Z) State[S, Z] {
	return NewState(func(state S) (Z, S) {
		val, nextState := s.Run(state)
		return f(val), nextState
	})
}

// AndThen returns a new State that applies the mapping function (returning another State) to the result.
func (s State[S, T]) AndThen[Z any](f func(T) State[S, Z]) State[S, Z] {
	return NewState(func(state S) (Z, S) {
		val, nextState := s.Run(state)
		return f(val).Run(nextState)
	})
}

// Get returns a State that produces the current state as the result.
func Get[S any]() State[S, S] {
	return NewState(func(s S) (S, S) {
		return s, s
	})
}

// Put returns a State that sets the state to newState.
func Put[S any](newState S) State[S, any] {
	return NewState(func(S) (any, S) {
		return nil, newState
	})
}

// Modify returns a State that modifies the state using the given function.
func Modify[S any](f func(S) S) State[S, any] {
	return NewState(func(s S) (any, S) {
		return nil, f(s)
	})
}
