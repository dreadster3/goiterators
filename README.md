# Go Iterators

A powerful and flexible iterator library for Go, providing both synchronous and asynchronous processing capabilities with comprehensive error handling.

## Features

- **Generic Iterator Interface**: Type-safe iterators using Go generics
- **Synchronous Algorithms**: Map, Filter, Take operations
- **Asynchronous Processing**: Parallel execution with MapAsync, FilterAsync, FlatMapAsync
- **Error Propagation**: Comprehensive error handling throughout iterator chains
- **Channel-based Async**: Efficient async processing using Go channels and goroutines
- **Easy Integration**: Works seamlessly with Go's standard library

## Installation

```bash
go get github.com/dreadster3/goiterators
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/dreadster3/goiterators"
)

func main() {
    // Create an iterator from a slice
    data := []int{1, 2, 3, 4, 5}
    iter := goiterators.NewIteratorFromSlice(data)
    
    // Chain operations
    result := goiterators.Map(
        goiterators.Filter(iter, func(x int) bool { return x%2 == 0 }),
        func(x int) int { return x * 2 },
    )
    
    // Collect results
    for item := range result.Next {
        fmt.Println(item) // Prints: 4, 8
    }
}
```

## API Reference

### Core Types

#### Iterator Interface

```go
type Iterator[T any] interface {
    Next(yield func(T) bool)        // Iterate over values
    INext(yield func(int, T) bool)  // Iterate with indices
    Err() error                     // Get any error that occurred
}
```

#### Result Type (Async)

```go
type Result[T any] struct {
    Value T
    Err   error
}
```

### Constructor Functions

- `NewIterator[T](iter.Seq2[int, T]) Iterator[T]` - Create from Go's standard iterator
- `NewIteratorErr[T](iter.Seq2[T, error]) Iterator[T]` - Create with error handling
- `NewIteratorFromSlice[T]([]T) Iterator[T]` - Create from slice
- `NewAsyncIterator[T](<-chan T) Iterator[T]` - Create async iterator from channel
- `NewAsyncIteratorErr[T](<-chan Result[T]) Iterator[T]` - Create async iterator with errors

### Synchronous Algorithms

#### Map
Transform each element using a function.

```go
func Map[T, U any](iter Iterator[T], fn func(T) U) Iterator[U]
```

#### Filter
Keep only elements that satisfy a predicate.

```go
func Filter[T any](iter Iterator[T], fn func(T) bool) Iterator[T]
```

#### Take
Take at most n elements from the iterator.

```go
func Take[T any](iter Iterator[T], n int) Iterator[T]
```

### Asynchronous Algorithms

#### MapAsync
Transform elements in parallel using goroutines.

```go
func MapAsync[T, U any](iter Iterator[T], fn func(T) U) Iterator[U]
```

#### FilterAsync
Filter elements in parallel.

```go
func FilterAsync[T any](iter Iterator[T], fn func(T) bool) Iterator[T]
```

#### FlatMapAsync
Transform each element into multiple results in parallel.

```go
func FlatMapAsync[T, U any](iter Iterator[T], fn func(T) []U) Iterator[U]
```

## Examples

### Basic Usage

```go
// Create iterator from slice
numbers := []int{1, 2, 3, 4, 5}
iter := goiterators.NewIteratorFromSlice(numbers)

// Transform with Map
doubled := goiterators.Map(iter, func(x int) int { return x * 2 })

// Collect results
var result []int
for item := range doubled.Next {
    result = append(result, item)
}
// result: [2, 4, 6, 8, 10]
```

### Chaining Operations

```go
data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
iter := goiterators.NewIteratorFromSlice(data)

// Chain: filter evens, take first 3, double each
result := goiterators.Map(
    goiterators.Take(
        goiterators.Filter(iter, func(x int) bool { return x%2 == 0 }),
        3,
    ),
    func(x int) int { return x * 2 },
)

// result: [4, 8, 12] (from 2, 4, 6 -> doubled)
```

### Async Processing

```go
data := []int{1, 2, 3, 4, 5}
iter := goiterators.NewIteratorFromSlice(data)

// Process in parallel
asyncResult := goiterators.MapAsync(iter, func(x int) int {
    time.Sleep(100 * time.Millisecond) // Simulate work
    return x * x
})

// Collect results (order may vary due to parallel execution)
var squares []int
for item := range asyncResult.Next {
    squares = append(squares, item)
}
```

### Error Handling

```go
// Create iterator that may produce errors
next := func(yield func(int, error) bool) {
    for i := 1; i <= 5; i++ {
        var err error
        if i == 3 {
            err = errors.New("error at 3")
        }
        if !yield(i, err) {
            return
        }
    }
}

iter := goiterators.NewIteratorErr(next)
mapped := goiterators.Map(iter, func(x int) int { return x * 2 })

// Process until error
for item := range mapped.Next {
    fmt.Println(item) // Prints: 2, 4
}

if err := mapped.Err(); err != nil {
    fmt.Println("Error:", err) // Prints: Error: error at 3
}
```

## Performance Considerations

### Synchronous vs Asynchronous

- **Synchronous**: Lower overhead, predictable execution order, better for CPU-bound tasks
- **Asynchronous**: Higher throughput for I/O-bound tasks, parallel execution, results may arrive out of order

### Memory Usage

- Iterators are lazy and process items on-demand
- Async iterators use channels for communication between goroutines
- Consider using `Take()` to limit processing when working with large datasets

## Testing

Run the comprehensive test suite:

```bash
go test ./...
```

The library includes extensive tests covering:
- Basic functionality
- Error propagation
- Async processing and parallelism
- Edge cases and error conditions

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.