package goiterators

import (
	"context"
	"iter"
	"slices"
	"sync"
)

// processAsync provides async processing with context cancellation support
// The worker function is called for each item with its index and can send zero or more results to the channel
func processAsync[T, U any](ctx context.Context, iter Iterator[T], worker func(context.Context, int, T, chan<- Result[U])) Iterator[U] {
	channel := make(chan Result[U])

	go func() {
		defer close(channel)
		wg := sync.WaitGroup{}

		for idx, item := range iter.INext {
			// Check for context cancellation
			select {
			case <-ctx.Done():
				channel <- Result[U]{Value: *new(U), Err: ctx.Err()}
				wg.Wait() // Wait for any pending goroutines
				return
			default:
			}

			// Check for error from underlying iterator
			if iter.Err() != nil {
				channel <- Result[U]{Value: *new(U), Err: iter.Err()}
				wg.Wait() // Wait for any pending goroutines
				return
			}

			wg.Add(1)
			go func(idx int, item T) {
				defer wg.Done()
				select {
				case <-ctx.Done():
					channel <- Result[U]{Value: *new(U), Err: ctx.Err()}
				default:
					worker(ctx, idx, item, channel)

				}
			}(idx, item)
		}

		// Wait for completion or context cancellation
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// All work completed normally
		case <-ctx.Done():
			// Context cancelled while waiting
			channel <- Result[U]{Value: *new(U), Err: ctx.Err()}
			wg.Wait() // Still wait for goroutines to finish
			return
		}

		// Final check for errors after processing all items
		if iter.Err() != nil {
			channel <- Result[U]{Value: *new(U), Err: iter.Err()}
		}
	}()

	return NewAsyncIteratorErr(channel)
}

// IMapAsyncCtx transforms each item using the provided function with index in parallel with context cancellation
func IMapAsyncCtx[T any, U any](ctx context.Context, iter Iterator[T], fn func(context.Context, int, T) (U, error)) Iterator[U] {
	return processAsync(ctx, iter, func(ctx context.Context, idx int, item T, ch chan<- Result[U]) {
		result, err := fn(ctx, idx, item)
		ch <- Result[U]{Value: result, Err: err}
	})
}

// MapAsync transforms each item using the provided function in parallel
func MapAsync[T any, U any](iter Iterator[T], fn func(T) U) Iterator[U] {
	return IMapAsync(iter, func(_ int, item T) U {
		return fn(item)
	})
}

// IMapAsync transforms each item using the provided function with index in parallel
func IMapAsync[T any, U any](iter Iterator[T], fn func(int, T) U) Iterator[U] {
	return IMapAsyncCtx(context.Background(), iter, func(ctx context.Context, i int, t T) (U, error) {
		return fn(i, t), nil
	})
}

// MapAsyncCtx transforms each item using the provided function in parallel with context cancellation
func MapAsyncCtx[T any, U any](ctx context.Context, iter Iterator[T], fn func(context.Context, T) (U, error)) Iterator[U] {
	return IMapAsyncCtx(ctx, iter, func(ctx context.Context, _ int, item T) (U, error) {
		return fn(ctx, item)
	})
}

// IFilterAsyncCtx returns only items that satisfy the predicate function with index in parallel with context cancellation
func IFilterAsyncCtx[T any](ctx context.Context, iter Iterator[T], fn func(context.Context, int, T) (bool, error)) Iterator[T] {
	return processAsync(ctx, iter, func(ctx context.Context, idx int, item T, ch chan<- Result[T]) {
		match, err := fn(ctx, idx, item)
		if err != nil {
			ch <- Result[T]{Value: *new(T), Err: err}
		} else if match {
			ch <- Result[T]{Value: item, Err: nil}
		}
	})
}

// FilterAsyncCtx returns only items that satisfy the predicate function in parallel with context cancellation
func FilterAsyncCtx[T any](ctx context.Context, iter Iterator[T], fn func(context.Context, T) (bool, error)) Iterator[T] {
	return IFilterAsyncCtx(ctx, iter, func(ctx context.Context, i int, t T) (bool, error) {
		return fn(ctx, t)
	})
}

// FilterAsync returns only items that satisfy the predicate function in parallel
func FilterAsync[T any](iter Iterator[T], fn func(T) bool) Iterator[T] {
	return IFilterAsync(iter, func(i int, t T) bool {
		return fn(t)
	})
}

// IFilterAsync returns only items that satisfy the predicate function with index in parallel
func IFilterAsync[T any](iter Iterator[T], fn func(int, T) bool) Iterator[T] {
	return IFilterAsyncCtx(context.Background(), iter, func(ctx context.Context, i int, t T) (bool, error) {
		return fn(i, t), nil
	})
}

// IFlatMapAsyncCtx transforms each item into multiple results with index in parallel with context cancellation
func IFlatMapAsyncCtx[T, U any](ctx context.Context, iter Iterator[T], fn func(context.Context, int, T) (iter.Seq[U], error)) Iterator[U] {
	return processAsync(ctx, iter, func(ctx context.Context, idx int, item T, ch chan<- Result[U]) {
		results, err := fn(ctx, idx, item)
		if err != nil {
			ch <- Result[U]{Value: *new(U), Err: err}
		} else {
			for result := range results {
				select {
				case <-ctx.Done():
					ch <- Result[U]{Value: *new(U), Err: ctx.Err()}
					return
				case ch <- Result[U]{Value: result, Err: nil}:
				}
			}
		}
	})
}

// FlatMapAsync transforms each item into multiple results in parallel
func FlatMapAsync[T, U any](iterator Iterator[T], fn func(T) iter.Seq[U]) Iterator[U] {
	return IFlatMapAsync(iterator, func(_ int, item T) iter.Seq[U] {
		return fn(item)
	})
}

// IFlatMapAsync transforms each item into multiple results with index in parallel
func IFlatMapAsync[T, U any](iterator Iterator[T], fn func(int, T) iter.Seq[U]) Iterator[U] {
	return IFlatMapAsyncCtx(context.Background(), iterator, func(ctx context.Context, i int, t T) (iter.Seq[U], error) {
		return fn(i, t), nil
	})
}

// FlatMapAsyncCtx transforms each item into multiple results in parallel with context cancellation
func FlatMapAsyncCtx[T, U any](ctx context.Context, iterator Iterator[T], fn func(context.Context, T) (iter.Seq[U], error)) Iterator[U] {
	return IFlatMapAsyncCtx(ctx, iterator, func(ctx context.Context, i int, t T) (iter.Seq[U], error) {
		return fn(ctx, t)
	})
}

// IForEachAsyncCtx applies the function to each item with index in parallel with context cancellation
func IForEachAsyncCtx[T any](ctx context.Context, iter Iterator[T], fn func(context.Context, int, T) error) error {
	processIterator := processAsync(ctx, iter, func(ctx context.Context, i int, t T, c chan<- Result[struct{}]) {
		c <- Result[struct{}]{Value: struct{}{}, Err: fn(ctx, i, t)}
	})

	_ = slices.Collect(processIterator.Next)

	return processIterator.Err()
}

// IForEachAsync applies the function to each item with index in parallel
func IForEachAsync[T any](iter Iterator[T], fn func(int, T) error) error {
	return IForEachAsyncCtx(context.Background(), iter, func(_ context.Context, i int, t T) error {
		return fn(i, t)
	})
}

// ForEachAsyncCtx applies the function to each item with index in parallel with context cancellation
func ForEachAsyncCtx[T any](ctx context.Context, iter Iterator[T], fn func(context.Context, T) error) error {
	return IForEachAsyncCtx(ctx, iter, func(ctx context.Context, i int, t T) error {
		return fn(ctx, t)
	})
}

// ForEachAsync applies the function to each item with index in parallel
func ForEachAsync[T any](iter Iterator[T], fn func(T) error) error {
	return ForEachAsyncCtx(context.Background(), iter, func(_ context.Context, t T) error {
		return fn(t)
	})
}
