# Mixed Operations

This example demonstrates seamlessly mixing synchronous and asynchronous operations.

## What it shows

- **Sync → Async → Sync** chains
- Performance optimization by using sync for cheap operations
- Same iterator interface works with both sync and async functions
- Strategic placement of operations for best performance

## Run

```bash
cd examples/mixed-operations
go run main.go
```

## Key concepts

- Any iterator can use both sync and async functions
- Use sync for fast operations (filtering, taking)
- Use async for expensive operations (heavy computation, I/O)
- Filter before expensive operations to reduce work
- Order operations strategically for optimal performance