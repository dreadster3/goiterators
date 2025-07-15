# Asynchronous Algorithms

This example demonstrates parallel processing with async algorithms.

## What it shows

- `MapAsync()` - Transform elements in parallel
- `FilterAsync()` - Filter elements in parallel
- `FlatMapAsync()` - Expand elements in parallel
- Performance benefits of parallel execution
- Results may arrive out of order

## Run

```bash
cd examples/async-algorithms
go run main.go
```

## Key concepts

- Each element is processed in a separate goroutine
- Significant performance improvement for expensive operations
- Results must be sorted if order matters
- Best for I/O-bound or CPU-intensive transformations