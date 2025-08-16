package promised

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
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
	result   PromiseResult[T]
	done     chan struct{}
	finished atomic.Bool
}

func New[T any](constrRoutine PromiseConstructor[T]) *Promise[T] {
	var once sync.Once

	promise := &Promise[T]{
		done: make(chan struct{}),
	}

	resolver := func(val T, err error) {
		promise.result = PromiseResult[T]{Value: val, Error: err}
		promise.finished.Store(true)
		close(promise.done)
	}

	go constrRoutine(func(value T) {
		once.Do(func() {
			resolver(value, nil)
		})
	}, func(err error) {
		once.Do(func() {
			var zero T
			resolver(zero, err)
		})
	})
	return promise
}

func (p *Promise[T]) Await() (T, error) {
	<-p.done
	return p.result.Value, p.result.Error
}

func (p *Promise[T]) IsFinished() bool {
	return p.finished.Load()
}

func AwaitAll[T any](promises ...*Promise[T]) []PromiseResult[T] {
	if len(promises) == 0 {
		return make([]PromiseResult[T], 0)
	}

	results := make([]PromiseResult[T], len(promises))

	var wg sync.WaitGroup
	wg.Add(len(promises))

	for i, p := range promises {
		go func(i int, p *Promise[T]) {
			defer wg.Done()
			val, err := p.Await()
			results[i] = PromiseResult[T]{Value: val, Error: err}
		}(i, p)
	}

	wg.Wait()
	return results
}

func AwaitAny[T any](promises ...*Promise[T]) (T, error) {
	var zero T

	if len(promises) == 0 {
		return zero, errors.New("promised: No Promises in call to AwaitFirstSuccess")
	}

	channel := make(chan PromiseResult[T], len(promises))

	for _, p := range promises {
		go func(p *Promise[T]) {
			val, err := p.Await()
			channel <- PromiseResult[T]{Value: val, Error: err}
		}(p)
	}

	result := <-channel
	return result.Value, result.Error
}

func AwaitAnySuccess[T any](promises ...*Promise[T]) (T, error) {
	var zero T

	if len(promises) == 0 {
		return zero, errors.New("promised: No Promises in call to AwaitFirstSuccess")
	}

	channel := make(chan PromiseResult[T], len(promises))

	for _, p := range promises {
		go func(p *Promise[T]) {
			val, err := p.Await()
			channel <- PromiseResult[T]{Value: val, Error: err}
		}(p)
	}

	for range promises {
		result := <-channel
		if result.Error == nil {
			return result.Value, nil
		}
	}

	errs := make([]error, len(promises))
	for i, p := range promises {
		errs[i] = p.result.Error
	}
	return zero, errors.Join(errs...)
}

func TimeoutPromise[T any](milliSeconds int64) *Promise[T] {
	promise := New(func(resolve func(value T), reject func(err error)) {
		time.AfterFunc(time.Duration(milliSeconds)*time.Millisecond, func() {
			reject(errors.New("promised: Timed out"))
		})
	})
	return promise
}
