package goiterators_test

import (
	"errors"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/dreadster3/goiterators"
	"github.com/stretchr/testify/assert"
)

func TestAsyncIterator(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	channel := make(chan int)

	go func() {
		defer close(channel)
		for _, item := range data {
			channel <- item
		}
	}()

	iterator := goiterators.NewAsyncIterator(channel)
	result := slices.Collect(iterator.Next)

	assert.Equal(t, data, result)
	assert.NoError(t, iterator.Err())
}

func TestAsyncIteratorWithError(t *testing.T) {
	data := []int{1, 2, 3}
	channel := make(chan goiterators.Result[int])

	go func() {
		defer close(channel)
		for _, item := range data {
			var err error
			if item == 2 {
				err = errors.New("error at 2")
			}
			channel <- goiterators.Result[int]{Value: item, Err: err}
		}
	}()

	iterator := goiterators.NewAsyncIteratorErr(channel)
	result := slices.Collect(iterator.Next)

	assert.Error(t, iterator.Err())
	assert.Equal(t, []int{1}, result)
}

func TestAsyncIteratorEmpty(t *testing.T) {
	channel := make(chan int)
	close(channel)

	iterator := goiterators.NewAsyncIterator(channel)
	result := slices.Collect(iterator.Next)

	assert.Empty(t, result)
	assert.NoError(t, iterator.Err())
}

func TestAsyncIteratorINext(t *testing.T) {
	data := []int{1, 2, 3}
	channel := make(chan int)

	go func() {
		defer close(channel)
		for _, item := range data {
			channel <- item
		}
	}()

	iterator := goiterators.NewAsyncIterator(channel)

	var indices []int
	var values []int

	for i, v := range iterator.INext {
		indices = append(indices, i)
		values = append(values, v)
	}

	assert.Equal(t, []int{0, 1, 2}, indices)
	assert.Equal(t, data, values)
	assert.NoError(t, iterator.Err())
}

func TestAsyncIteratorConcurrentAccess(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	channel := make(chan int)

	go func() {
		defer close(channel)
		for _, item := range data {
			channel <- item
			time.Sleep(5 * time.Millisecond) // Simulate slow producer
		}
	}()

	iterator := goiterators.NewAsyncIterator(channel)

	// Test concurrent access from multiple goroutines
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allResults []int

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			var localResults []int
			for item := range iterator.Next {
				localResults = append(localResults, item)
				// Stop after getting some items to allow other goroutines to consume
				if len(localResults) >= 2 {
					break
				}
			}

			mu.Lock()
			allResults = append(allResults, localResults...)
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	assert.NoError(t, iterator.Err())
	assert.NotEmpty(t, allResults)
	// Should have gotten some results from concurrent consumption
	assert.LessOrEqual(t, len(allResults), len(data))
}

func TestAsyncIteratorErrorTiming(t *testing.T) {
	channel := make(chan goiterators.Result[int])

	go func() {
		defer close(channel)
		// Send some good values
		channel <- goiterators.Result[int]{Value: 1, Err: nil}
		channel <- goiterators.Result[int]{Value: 2, Err: nil}

		// Send an error
		channel <- goiterators.Result[int]{Value: 0, Err: errors.New("async error")}

		// Send more values after error (should not be consumed)
		channel <- goiterators.Result[int]{Value: 3, Err: nil}
		channel <- goiterators.Result[int]{Value: 4, Err: nil}
	}()

	iterator := goiterators.NewAsyncIteratorErr(channel)
	result := slices.Collect(iterator.Next)

	assert.Error(t, iterator.Err())
	assert.Equal(t, "async error", iterator.Err().Error())
	assert.Equal(t, []int{1, 2}, result)
}

func TestAsyncIteratorChannelClosure(t *testing.T) {
	channel := make(chan int)

	// Close channel immediately
	close(channel)

	iterator := goiterators.NewAsyncIterator(channel)

	// Should handle closed channel gracefully
	result := slices.Collect(iterator.Next)

	assert.Empty(t, result)
	assert.NoError(t, iterator.Err())
}

func TestAsyncIteratorINextWithError(t *testing.T) {
	channel := make(chan goiterators.Result[int])

	go func() {
		defer close(channel)
		channel <- goiterators.Result[int]{Value: 10, Err: nil}
		channel <- goiterators.Result[int]{Value: 20, Err: nil}
		channel <- goiterators.Result[int]{Value: 0, Err: errors.New("INext error")}
		channel <- goiterators.Result[int]{Value: 30, Err: nil} // Should not be processed
	}()

	iterator := goiterators.NewAsyncIteratorErr(channel)

	var indices []int
	var values []int

	for i, v := range iterator.INext {
		indices = append(indices, i)
		values = append(values, v)
	}

	assert.Error(t, iterator.Err())
	assert.Equal(t, "INext error", iterator.Err().Error())
	assert.Equal(t, []int{0, 1}, indices)
	assert.Equal(t, []int{10, 20}, values)
}
