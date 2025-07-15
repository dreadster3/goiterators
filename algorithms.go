package goiterators

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
