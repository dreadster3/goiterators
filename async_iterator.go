package goiterators

// Result wraps a value with an optional error for async operations
type Result[T any] struct {
	Value T
	Err   error
}

type asyncIterator[T any] struct {
	dataIn <-chan Result[T]
	err    error
}

// NewAsyncIterator creates an async iterator from a channel of values
func NewAsyncIterator[T any](dataIn <-chan T) Iterator[T] {
	channel := make(chan Result[T])
	go func() {
		defer close(channel)
		i := 0
		for item := range dataIn {
			channel <- Result[T]{
				Value: item,
				Err:   nil,
			}
			i += 1
		}
	}()

	return NewAsyncIteratorErr(channel)
}

// NewAsyncIteratorErr creates an async iterator from a channel of Results
func NewAsyncIteratorErr[T any](dataIn <-chan Result[T]) Iterator[T] {
	return &asyncIterator[T]{
		dataIn: dataIn,
	}
}

func (it *asyncIterator[T]) Next(yield func(T) bool) {
	for item := range it.dataIn {
		if item.Err != nil {
			it.err = item.Err
			return
		}

		if !yield(item.Value) {
			return
		}
	}
}

func (it *asyncIterator[T]) INext(yield func(int, T) bool) {
	i := 0
	for item := range it.dataIn {
		if item.Err != nil {
			it.err = item.Err
			return
		}

		if !yield(i, item.Value) {
			return
		}

		i += 1
	}
}

func (it *asyncIterator[T]) Err() error {
	return it.err
}
