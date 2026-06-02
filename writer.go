package main

// Writer represents a value of type T alongside a log (side-channel) of type []W.
type Writer[W, T any] struct {
	value T
	logs  []W
}

// NewWriter returns a new Writer with the given value and logs.
func NewWriter[W, T any](value T, logs ...W) Writer[W, T] {
	return Writer[W, T]{value: value, logs: logs}
}

// Value returns the value of the computation.
func (w Writer[W, T]) Value() T {
	return w.value
}

// Logs returns the accumulated logs.
func (w Writer[W, T]) Logs() []W {
	return w.logs
}

// Map returns a new Writer that applies the mapping function to the value.
func (w Writer[W, T]) Map[Z any](f func(T) Z) Writer[W, Z] {
	return NewWriter(f(w.value), w.logs...)
}

// AndThen returns a new Writer that applies the mapping function (returning another Writer) to the value,
// accumulating logs from both.
func (w Writer[W, T]) AndThen[Z any](f func(T) Writer[W, Z]) Writer[W, Z] {
	next := f(w.value)
	newLogs := make([]W, 0, len(w.logs)+len(next.logs))
	newLogs = append(newLogs, w.logs...)
	newLogs = append(newLogs, next.logs...)
	return NewWriter(next.value, newLogs...)
}

// Tell returns a Writer that has a nil value but contains the given log.
func Tell[W any](log W) Writer[W, any] {
	return NewWriter[W, any](nil, log)
}
