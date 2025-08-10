# Promised

A Go library that provides JavaScript-style Promise syntax for asynchronous operations, built on top of Go's native goroutines and channels.

## Overview

While Go's goroutines and channels are powerful primitives for concurrent programming, developers coming from JavaScript or other languages might find the Promise pattern more familiar and intuitive. This library bridges that gap by providing a Promise-like API that internally uses goroutines and channels.

## Installation

```bash
go get github.com/wAIfu-DEV/promised
```

## Quick Start

```go
package main

import (
    "errors"
    "fmt"
    "github.com/wAIfu-DEV/promised"
)

func main() {
    // Create a new Promise
    promise := promised.New(func(resolve func(value int), reject func(err error)) {
        // This function runs in a goroutine
        fmt.Println("Running async operation...")
        
        // Simulate some work
        if true {
            resolve(42) // Success case
        } else {
            reject(errors.New("something went wrong")) // Error case
        }
    })
    
    // Wait for the result
    value, err := promise.Await()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Result: %d\n", value)
    }
}
```

## API Reference

### Types

#### `PromiseConstructor[T any]`
```go
type PromiseConstructor[T any] func(
    resolve func(value T),
    reject func(err error),
)
```
A function type that defines the Promise constructor. It receives two callback functions:
- `resolve`: Call this function with a value to fulfill the Promise successfully
- `reject`: Call this function with an error to reject the Promise

#### `PromiseResult[T any]`
```go
type PromiseResult[T any] struct {
    Value T
    Error error
}
```
Internal structure that holds the result of a Promise operation.

#### `Promise[T any]`
```go
type Promise[T any] struct {
    ch <-chan PromiseResult[T]
}
```
The main Promise type that wraps a channel for result communication.

### Functions

#### `New[T any](constrRoutine PromiseConstructor[T]) *Promise[T]`

Creates a new Promise that will execute the provided constructor function in a goroutine.

**Parameters:**
- `constrRoutine`: A function that receives `resolve` and `reject` callbacks

**Returns:**
- A pointer to a new `Promise[T]`

**Example:**
```go
promise := promised.New(func(resolve func(value string), reject func(err error)) {
    // Your async logic here
    resolve("Hello, World!")
})
```

#### `(p *Promise[T]) Await() (T, error)`

Blocks until the Promise is resolved or rejected, then returns the result.

**Returns:**
- `T`: The resolved value (zero value if Promise was rejected)
- `error`: The rejection error (nil if Promise was resolved)

**Example:**
```go
value, err := promise.Await()
if err != nil {
    // Handle error
} else {
    // Use value
}
```

## Examples

### Basic Success Case

```go
promise := promised.New(func(resolve func(value int), reject func(err error)) {
    // Simulate some async work
    time.Sleep(100 * time.Millisecond)
    resolve(100)
})

result, err := promise.Await()
// result = 100, err = nil
```

### Error Handling

```go
promise := promised.New(func(resolve func(value int), reject func(err error)) {
    // Simulate an error condition
    reject(errors.New("operation failed"))
})

result, err := promise.Await()
// result = 0 (zero value), err = "operation failed"
```

### Working with Different Types

```go
// String Promise
stringPromise := promised.New(func(resolve func(value string), reject func(err error)) {
    resolve("Hello, Go!")
})

// Struct Promise
type User struct {
    ID   int
    Name string
}

userPromise := promised.New(func(resolve func(value User), reject func(err error)) {
    user := User{ID: 1, Name: "Alice"}
    resolve(user)
})

str, _ := stringPromise.Await()
user, _ := userPromise.Await()
```

## Comparison with Traditional Go Patterns

### Traditional Goroutine Approach
```go
func withGoroutine() int {
    channel := make(chan int)
    
    go func() {
        fmt.Println("Running Goroutine")
        channel <- 10
    }()
    
    value, ok := <-channel
    if !ok {
        fmt.Println("Failed to get value")
    }
    return value
}
```

### Promise Approach
```go
func withPromise() int {
    promise := promised.New(func(resolve func(value int), reject func(err error)) {
        fmt.Println("Running Promise")
        resolve(10)
    })
    
    value, err := promise.Await()
    if err != nil {
        fmt.Println(err.Error())
    }
    return value
}
```

## Features

- **Type Safe**: Full generic support for any type `T`
- **Goroutine Based**: Internally uses Go's native concurrency primitives
- **JavaScript-like Syntax**: Familiar Promise pattern for developers coming from JavaScript
- **Error Handling**: Built-in error propagation through the `reject` callback
- **Thread Safe**: Uses `sync.Once` to ensure resolve/reject can only be called once

## Implementation Notes

- Each Promise runs its constructor function in a separate goroutine
- The `resolve` and `reject` functions can only be called once (enforced by `sync.Once`)
- Calling `Await()` blocks until the Promise is settled (resolved or rejected)
- The underlying channel is automatically closed after the Promise settles
- Zero values are returned for the value type when a Promise is rejected

## License

This project is licensed under the terms specified in the repository.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.
