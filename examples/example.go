package _promised_example

import (
	"errors"
	"fmt"
	"testing"

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
