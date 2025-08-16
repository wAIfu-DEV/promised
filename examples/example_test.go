package _promised_example

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/wAIfu-DEV/promised"
)

// Idiomatic way to accomplish "Promise" like async in go
func withGoroutine() int {

	channel := make(chan int)

	go func() {
		fmt.Println("Running Goroutine")
		channel <- 10
	}()

	value, ok := <-channel
	if !ok {
		fmt.Println("Failed to get value with Goroutine.")
	}
	return value
}

// Using this package, syntax is more in line with JS, but internals are similar
func withPromise() int {

	promise := promised.New(func(resolve func(value int), reject func(err error)) {
		// This function will be run as a goroutine ASAP
		fmt.Println("Running Promise")

		if true {
			resolve(10) // Finishes call to Await with value.
		} else {
			// Finishes call to Await with error.
			reject(errors.New("failed to get value with Promise"))
		}
	})

	value, err := promise.Await() // Blocks until value or error is available
	if err != nil {
		fmt.Println(err.Error())
	}
	return value
}

func TestExamples(t *testing.T) {
	_ = withGoroutine()
	_ = withPromise()
}

func TestMultipleHandlers(t *testing.T) {
	fmt.Println("Testing promised.AwaitAll")

	promise1 := promised.New(func(resolve func(value int), reject func(err error)) {
		reject(errors.New("reject1"))
	})

	promise2 := promised.New(func(resolve func(value int), reject func(err error)) {
		resolve(1)
	})

	promise3 := promised.New(func(resolve func(value int), reject func(err error)) {
		reject(errors.New("reject3"))
	})

	results := promised.AwaitAll(promise1, promise2, promise3)

	for i, res := range results {
		if res.Error != nil {
			fmt.Printf("Promise %d: rejected.\n", i)
		} else {
			fmt.Printf("Promise %d: resolved with value: %d.\n", i, res.Value)
		}
	}

	fmt.Println("Testing promised.AwaitAnySuccess")

	promise4 := promised.New(func(resolve func(value int), reject func(err error)) {
		reject(errors.New("reject4"))
	})

	promise5 := promised.New(func(resolve func(value int), reject func(err error)) {
		time.Sleep(1500 * time.Millisecond)
		resolve(3)
	})

	value, err := promised.AwaitAnySuccess(promise4, promise5)
	if err != nil {
		fmt.Println("should have gotten value from call to promised.AwaitAnySuccess")
	} else {
		fmt.Printf("First Resolved Promise: resolved with value: %d\n", value)
	}

	fmt.Println("Testing promised.AwaitAny with Timeout")

	promise6 := promised.TimeoutPromise[int](1500)

	_, err = promised.AwaitAny(promise6)
	if err != nil {
		fmt.Printf("call to AwaitAny ended with expected error: %s\n", err.Error())
	}
}
