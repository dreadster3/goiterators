# Error Handling

This example demonstrates comprehensive error handling patterns.

## What it shows

- Creating iterators that can produce errors
- Error propagation through transformations
- Error handling in both sync and async operations
- Graceful error handling patterns
- Partial results before errors occur

## Run

```bash
cd examples/error-handling
go run main.go
```

## Key concepts

- `NewIteratorErr()` creates iterators that can handle errors
- Errors stop processing immediately
- `iter.Err()` checks for errors after iteration
- Both sync and async operations propagate errors correctly
- You can get partial results before the error occurred