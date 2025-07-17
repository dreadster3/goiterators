package goiterators

import "iter"

// Map transforms each item using the provided function
func Map[T any, U any](iter Iterator[T], fn func(T) U) Iterator[U] {
	return newIterator(func(self *iterator[U], yield func(int, U) bool) {
		for idx, item := range iter.INext {
			if !yield(idx, fn(item)) {
				return
			}
		}

		if iter.Err() != nil {
			self.err = iter.Err()
		}
	})
}

// IMap transforms each item using the provided function
func IMap[T, U any](iter Iterator[T], fn func(int, T) U) Iterator[U] {
	return newIterator(func(self *iterator[U], yield func(int, U) bool) {
		for idx, item := range iter.INext {
			if !yield(idx, fn(idx, item)) {
				return
			}
		}

		if iter.Err() != nil {
			self.err = iter.Err()
		}
	})
}

// Filter returns only items that satisfy the predicate function
func Filter[T any](iter Iterator[T], fn func(T) bool) Iterator[T] {
	return newIterator(func(self *iterator[T], yield func(int, T) bool) {
		for idx, item := range iter.INext {
			if fn(item) {
				if !yield(idx, item) {
					return
				}
			}
		}
		if iter.Err() != nil {
			self.err = iter.Err()
		}
	})
}

// IFilter returns only items that satisfy the predicate function with index
func IFilter[T any](iter Iterator[T], fn func(int, T) bool) Iterator[T] {
	return newIterator(func(self *iterator[T], yield func(int, T) bool) {
		for idx, item := range iter.INext {
			if fn(idx, item) {
				if !yield(idx, item) {
					return
				}
			}
		}
		if iter.Err() != nil {
			self.err = iter.Err()
		}
	})
}

// Take returns at most n items from the iterator
func Take[T any](iter Iterator[T], n int) Iterator[T] {
	return newIterator(func(self *iterator[T], yield func(int, T) bool) {
		if n <= 0 {
			return
		}

		for idx, item := range iter.INext {
			if !yield(idx, item) {
				return
			}

			if idx >= (n - 1) {
				return
			}
		}

		if iter.Err() != nil {
			self.err = iter.Err()
		}
	})
}

// FlatMap transforms each item into multiple results using iter.Seq
func FlatMap[T, U any](iter Iterator[T], fn func(T) iter.Seq[U]) Iterator[U] {
	return newIterator(func(self *iterator[U], yield func(int, U) bool) {
		outputIdx := 0
		for _, item := range iter.INext {
			for result := range fn(item) {
				if !yield(outputIdx, result) {
					return
				}
				outputIdx++
			}
		}

		if iter.Err() != nil {
			self.err = iter.Err()
		}
	})
}
