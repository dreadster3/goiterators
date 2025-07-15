package goiterators

import (
	"iter"
	"slices"
)

// Iterator provides sequential access to items with optional error handling
type Iterator[T any] interface {
	Next(yield func(T) bool)
	INext(yield func(int, T) bool)
	Err() error
}

type nextFunc[T any] func(self *iterator[T], yield func(int, T) bool)

type iterator[T any] struct {
	next nextFunc[T]
	err  error
}

// newIterator creates an iterator with error checking wrapper
func newIterator[T any](next nextFunc[T]) *iterator[T] {
	return &iterator[T]{
		next: func(self *iterator[T], yield func(int, T) bool) {
			if self.err != nil {
				return
			}
			next(self, yield)
		},
	}
}

// NewIterator creates an iterator from a standard Go iter.Seq2[int, T]
func NewIterator[T any](next iter.Seq2[int, T]) Iterator[T] {
	return &iterator[T]{
		next: func(self *iterator[T], yield func(int, T) bool) {
			for i, item := range next {
				if !yield(i, item) {
					return
				}
			}
		},
	}
}

// NewIteratorErr creates an iterator that handles errors from iter.Seq2[T, error]
func NewIteratorErr[T any](next iter.Seq2[T, error]) Iterator[T] {
	return &iterator[T]{
		next: func(self *iterator[T], yield func(int, T) bool) {
			if self.Err() != nil {
				return
			}

			i := 0
			next(func(item T, err error) bool {
				if err != nil {
					self.err = err
					return false
				}

				result := yield(i, item)
				i += 1
				return result
			})
		},
	}
}

// NewIteratorFromSlice creates an iterator from a slice
func NewIteratorFromSlice[T any](slice []T) Iterator[T] {
	return NewIterator(slices.All(slice))
}

func (it *iterator[T]) Err() error {
	return it.err
}

func (it *iterator[T]) Next(yield func(T) bool) {
	it.next(it, func(_ int, item T) bool {
		return yield(item)
	})
}

func (it *iterator[T]) INext(yield func(int, T) bool) {
	it.next(it, yield)
}
