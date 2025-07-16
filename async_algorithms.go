package goiterators

import "sync"

// processAsyncI provides async processing with index access
// The worker function is called for each item with its index and can send zero or more results to the channel
func processAsync[T, U any](iter Iterator[T], worker func(int, T, chan<- Result[U])) Iterator[U] {
	channel := make(chan Result[U])

	go func() {
		defer close(channel)
		wg := sync.WaitGroup{}

		for idx, item := range iter.INext {
			// Check for error from underlying iterator
			if iter.Err() != nil {
				channel <- Result[U]{Value: *new(U), Err: iter.Err()}
				return
			}

			wg.Add(1)
			go func(idx int, item T) {
				defer wg.Done()
				worker(idx, item, channel)
			}(idx, item)
		}
		wg.Wait()

		// Final check for errors after processing all items
		if iter.Err() != nil {
			channel <- Result[U]{Value: *new(U), Err: iter.Err()}
		}
	}()

	return NewAsyncIteratorErr(channel)
}

// MapAsync transforms each item using the provided function in parallel
func MapAsync[T any, U any](iter Iterator[T], fn func(T) U) Iterator[U] {
	return processAsync(iter, func(_ int, item T, ch chan<- Result[U]) {
		ch <- Result[U]{Value: fn(item), Err: nil}
	})
}

// IMapAsync transforms each item using the provided function with index in parallel
func IMapAsync[T any, U any](iter Iterator[T], fn func(int, T) U) Iterator[U] {
	return processAsync(iter, func(idx int, item T, ch chan<- Result[U]) {
		ch <- Result[U]{Value: fn(idx, item), Err: nil}
	})
}

// FilterAsync returns only items that satisfy the predicate function in parallel
func FilterAsync[T any](iter Iterator[T], fn func(T) bool) Iterator[T] {
	return processAsync(iter, func(_ int, item T, ch chan<- Result[T]) {
		if fn(item) {
			ch <- Result[T]{Value: item, Err: nil}
		}
	})
}

// IFilterAsync returns only items that satisfy the predicate function with index in parallel
func IFilterAsync[T any](iter Iterator[T], fn func(int, T) bool) Iterator[T] {
	return processAsync(iter, func(idx int, item T, ch chan<- Result[T]) {
		if fn(idx, item) {
			ch <- Result[T]{Value: item, Err: nil}
		}
	})
}

// FlatMapAsync transforms each item into multiple results in parallel
func FlatMapAsync[T, U any](iter Iterator[T], fn func(T) []U) Iterator[U] {
	return processAsync(iter, func(_ int, item T, ch chan<- Result[U]) {
		for _, result := range fn(item) {
			ch <- Result[U]{Value: result, Err: nil}
		}
	})
}

// IFlatMapAsync transforms each item into multiple results with index in parallel
func IFlatMapAsync[T, U any](iter Iterator[T], fn func(int, T) []U) Iterator[U] {
	return processAsync(iter, func(idx int, item T, ch chan<- Result[U]) {
		for _, result := range fn(idx, item) {
			ch <- Result[U]{Value: result, Err: nil}
		}
	})
}
