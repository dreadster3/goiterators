# Indexed Operations

This example demonstrates the indexed versions of all algorithms that provide access to element indices.

## What it shows

- `IMap()` - Transform with access to both index and value
- `IFilter()` - Filter based on index and value conditions
- `IMapAsync()`, `IFilterAsync()`, `IFlatMapAsync()` - Async versions with index access
- Chaining indexed and regular operations

## Run

```bash
cd examples/indexed-operations
go run main.go
```

## Key concepts

- All algorithms have indexed (`I`-prefixed) versions
- Index-aware functions receive `(int, T)` instead of just `T`
- Useful for position-dependent operations
- Can be mixed with regular (non-indexed) operations
- Async indexed operations maintain parallel processing benefits