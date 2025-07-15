# Basic Usage

This example demonstrates the fundamental operations of the iterator library.

## What it shows

- Creating iterators from slices
- Using `Next()` to iterate over values
- Using `INext()` to iterate with indices
- Collecting results into slices

## Run

```bash
cd examples/basic-usage
go run main.go
```

## Key concepts

- `NewIteratorFromSlice()` creates an iterator from any slice
- `slices.Collect()` gathers all iterator results into a slice
- Iterators support both value-only and indexed iteration