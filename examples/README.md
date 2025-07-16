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

# Indexed operations
cd examples/indexed-operations && go run main.go

# Context cancellation
cd examples/context-cancellation && go run main.go

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

### 🔢 [indexed-operations/](indexed-operations/)
Indexed versions of all algorithms with access to element positions.

### 🚫 [context-cancellation/](context-cancellation/)
Context cancellation and timeout support for async operations.

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
goiterators.FlatMap(iter, func(x int) []int { return []int{x, x*2} })

// Transform with index (sync)
goiterators.IMap(iter, func(idx int, x int) int { return x * idx })
goiterators.IFilter(iter, func(idx int, x int) bool { return idx%2 == 0 })

// Transform (async)
goiterators.MapAsync(iter, expensiveFunc)
goiterators.FilterAsync(iter, expensivePredicate)
goiterators.FlatMapAsync(iter, expandFunc)

// Transform with index (async)
goiterators.IMapAsync(iter, expensiveFuncWithIndex)
goiterators.IFilterAsync(iter, expensivePredicateWithIndex)
goiterators.IFlatMapAsync(iter, expandFuncWithIndex)

// Collect results
result := slices.Collect(iter.Next)

// Check for errors
if err := iter.Err(); err != nil { /* handle */ }
```