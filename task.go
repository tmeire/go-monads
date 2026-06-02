package main

// Task represents a lazy computation that produces a Result[T].
type Task[T any] struct {
	run func() Result[T]
}

// NewTask returns a new Task with the given computation.
func NewTask[T any](f func() Result[T]) Task[T] {
	return Task[T]{run: f}
}

// Run executes the Task synchronously and returns the Result.
func (t Task[T]) Run() Result[T] {
	return t.run()
}

// RunAsync executes the Task in a new goroutine and returns a channel that will receive the Result.
func (t Task[T]) RunAsync() <-chan Result[T] {
	ch := make(chan Result[T], 1)
	go func() {
		ch <- t.run()
	}()
	return ch
}

// Map returns a new Task that applies the mapping function to the success value of this Task.
func (t Task[T]) Map[Z any](f func(T) Z) Task[Z] {
	return NewTask(func() Result[Z] {
		res := t.run()
		return res.Map(f)
	})
}

// AndThen returns a new Task that applies the mapping function (returning another Task) to the success value of this Task.
func (t Task[T]) AndThen[Z any](f func(T) Task[Z]) Task[Z] {
	return NewTask(func() Result[Z] {
		res := t.run()
		if res.IsErr() {
			_, err := res.Get()
			return Err[Z](err)
		}
		val, _ := res.Get()
		return f(val).run()
	})
}
