# ForEach Examples

This example demonstrates the usage of ForEach functions in the goiterators library.

## Running the Example

```bash
cd examples/foreach
go run main.go
```

## What This Example Shows

### 1. Basic ForEach
- Apply a function to each element without returning values
- Useful for side effects like printing, logging, or updating external state

### 2. IForEach (Indexed ForEach)
- Same as ForEach but provides the index of each element
- Useful when you need to know the position of elements

### 3. ForEach with Side Effects
- Demonstrates how to use ForEach to calculate aggregations
- Shows how to modify external variables from within the ForEach function

### 4. ForEachAsync
- Parallel processing of elements using goroutines
- Significantly faster for I/O-bound operations
- Results may be processed out of order
- Uses the actual `ForEachAsync` function from the library

### 5. IForEachAsync (Indexed Async ForEach)
- Parallel processing with index information
- Combines the benefits of async processing with index awareness
- Uses the actual `IForEachAsync` function from the library

### 6. ForEach with Chained Operations
- Shows how ForEach works at the end of a processing chain
- Demonstrates the composition of filter, take, and forEach operations

## Key Points

- **ForEach vs Map**: ForEach is for side effects, Map is for transformations
- **Sync vs Async**: Async versions are better for I/O-bound tasks
- **Error Handling**: All ForEach functions return errors for proper error handling
- **Performance**: Async versions process elements in parallel, which can be much faster for certain workloads

## Expected Output

The example will show:
1. Sequential processing of elements
2. Indexed processing with element positions
3. Sum calculation using side effects
4. Parallel processing with timing information
5. Chained operations results

The async examples will demonstrate parallel execution by showing that multiple elements are processed simultaneously, reducing total execution time.