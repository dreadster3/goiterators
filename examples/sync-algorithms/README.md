# Synchronous Algorithms

This example demonstrates the core synchronous transformation algorithms.

## What it shows

- `Map()` - Transform each element
- `Filter()` - Keep elements matching a condition
- `Take()` - Limit the number of elements
- `FlatMap()` - Transform each element into multiple results
- Chaining operations together
- Working with different data types

## Run

```bash
cd examples/sync-algorithms
go run main.go
```

## Key concepts

- Operations execute sequentially in the order they appear
- Chaining creates readable functional-style code
- Each algorithm preserves the iterator interface
- Works with any data type using Go generics