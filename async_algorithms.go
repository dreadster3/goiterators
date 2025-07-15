package goiterators

import "sync"

// processAsync provides a common pattern for async processing of iterators
// The worker function is called for each item and can send zero or more results to the channel
func processAsync[T, U any](iter Iterator[T], worker func(T, chan<- Result[U])) Iterator[U] {
	channel := make(chan Result[U])

	go func() {
		defer close(channel)
		wg := sync.WaitGroup{}

		for item := range iter.Next {
			// Check for error from underlying iterator
			if iter.Err() != nil {
				channel <- Result[U]{Value: *new(U), Err: iter.Err()}
				return
			}

			wg.Add(1)
			go func(item T) {
				defer wg.Done()
				worker(item, channel)
			}(item)
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
	return processAsync(iter, func(item T, ch chan<- Result[U]) {
		ch <- Result[U]{Value: fn(item), Err: nil}
	})
}

// FilterAsync returns only items that satisfy the predicate function in parallel
func FilterAsync[T any](iter Iterator[T], fn func(T) bool) Iterator[T] {
	return processAsync(iter, func(item T, ch chan<- Result[T]) {
		if fn(item) {
			ch <- Result[T]{Value: item, Err: nil}
		}
	})
}

// FlatMapAsync transforms each item into multiple results in parallel
func FlatMapAsync[T, U any](iter Iterator[T], fn func(T) []U) Iterator[U] {
	return processAsync(iter, func(item T, ch chan<- Result[U]) {
		for _, result := range fn(item) {
			ch <- Result[U]{Value: result, Err: nil}
		}
	})
}
