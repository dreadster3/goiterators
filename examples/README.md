# Examples

This directory contains focused examples demonstrating the usage of the Go Iterators library. Each example is in its own folder with a clear, concise demonstration.

## Running Examples

```bash
# Basic iterator operations
cd examples/basic-usage && go run main.go

# Synchronous algorithms
cd examples/sync-algorithms && go run main.go

# Asynchronous algorithms  
cd examples/async-algorithms && go run main.go

# Mixed sync/async operations
cd examples/mixed-operations && go run main.go

# Error handling patterns
cd examples/error-handling && go run main.go
```

## Examples Overview

### 🚀 [basic-usage/](basic-usage/)
Fundamental iterator operations - creating, iterating, and collecting results.

### 🔄 [sync-algorithms/](sync-algorithms/) 
Synchronous transformations with `Map`, `Filter`, `Take`, and chaining.

### ⚡ [async-algorithms/](async-algorithms/)
Parallel processing with `MapAsync`, `FilterAsync`, and `FlatMapAsync`.

### 🔀 [mixed-operations/](mixed-operations/)
**Key Feature**: Seamlessly mixing sync and async operations in the same chain.

### ⚠️ [error-handling/](error-handling/)
Comprehensive error propagation and graceful handling patterns.

## Key Insight

The library's main strength is **unified iterator interface** - any iterator can use both sync and async functions:

```go
iter := goiterators.NewIteratorFromSlice(data)

// Mix and match as needed
result := goiterators.Take(                    // Sync
    goiterators.MapAsync(                      // Async  
        goiterators.Filter(iter, predicate),   // Sync
        expensiveTransform,
    ),
    10,
)
```

## Quick Reference

```go
// Create iterators
iter := goiterators.NewIteratorFromSlice([]int{1, 2, 3})
iter := goiterators.NewIteratorErr(errorFunction)

// Transform (sync)
goiterators.Map(iter, func(x int) int { return x * 2 })
goiterators.Filter(iter, func(x int) bool { return x > 0 })
goiterators.Take(iter, 5)

// Transform (async)
goiterators.MapAsync(iter, expensiveFunc)
goiterators.FilterAsync(iter, expensivePredicate)
goiterators.FlatMapAsync(iter, expandFunc)

// Collect results
result := slices.Collect(iter.Next)

// Check for errors
if err := iter.Err(); err != nil { /* handle */ }
```