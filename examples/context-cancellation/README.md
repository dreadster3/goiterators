# Context Cancellation

This example demonstrates context cancellation support for async algorithms, enabling graceful shutdown and timeout handling.

## What it shows

- **Context-aware async functions**: `MapAsyncCtx`, `FilterAsyncCtx`, `FlatMapAsyncCtx` and indexed versions
- **Cancellation**: Stop processing early with `context.WithCancel`
- **Timeouts**: Automatic cancellation with `context.WithTimeout`
- **Proper cleanup**: Goroutines respect context cancellation
- **Error handling**: Context errors propagate through iterator chains
- **Mixed operations**: Combine context-aware and regular operations

## Run

```bash
cd examples/context-cancellation
go run main.go
```

## Key concepts

- Context-aware functions take `context.Context` as first parameter
- Work functions can return errors (including context cancellation)
- Processing stops immediately on context cancellation/timeout
- Partial results are available before cancellation
- Context cancellation is checked both between items and within work functions
- Can be mixed with regular (non-context) operations seamlessly