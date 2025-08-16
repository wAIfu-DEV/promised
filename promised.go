package promised

import (
	"sync"
)

type PromiseConstructor[T any] func(
	resolve func(value T),
	reject func(err error),
)

type PromiseResult[T any] struct {
	Value T
	Error error
}

type Promise[T any] struct {
	ch <-chan PromiseResult[T]
}

func New[T any](constrRoutine PromiseConstructor[T]) *Promise[T] {
	var channel chan PromiseResult[T] = make(chan PromiseResult[T], 1)
	var once sync.Once

	go constrRoutine(func(value T) {
		once.Do(func() {
			channel <- PromiseResult[T]{Value: value, Error: nil}
			close(channel)
		})
	}, func(err error) {
		once.Do(func() {
			channel <- PromiseResult[T]{Error: err}
			close(channel)
		})
	})

	return &Promise[T]{ch: channel}
}

func (p *Promise[T]) Await() (T, error) {
	result := <-p.ch
	return result.Value, result.Error
}

func (p *Promise[T]) IsFinished() bool {
	return len(p.ch) > 0
}
