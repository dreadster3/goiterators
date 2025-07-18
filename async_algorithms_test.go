package goiterators_test

import (
	"errors"
	"fmt"
	"iter"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dreadster3/goiterators"
	"github.com/stretchr/testify/assert"
)

func TestMapAsync(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iterator := goiterators.NewIteratorFromSlice(data)

	mapped := goiterators.MapAsync(iterator, func(item int) int {
		time.Sleep(10 * time.Millisecond)
		return item * 2
	})

	result := slices.Collect(mapped.Next)
	slices.Sort(result)

	expected := []int{2, 4, 6, 8, 10}
	assert.Equal(t, expected, result)
	assert.NoError(t, mapped.Err())
}

func TestMapAsyncParallelism(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iterator := goiterators.NewIteratorFromSlice(data)

	// Track concurrent executions
	var concurrent int64
	var maxConcurrent int64

	start := time.Now()
	mapped := goiterators.MapAsync(iterator, func(item int) int {
		// Increment concurrent counter
		current := atomic.AddInt64(&concurrent, 1)
		if current > atomic.LoadInt64(&maxConcurrent) {
			atomic.StoreInt64(&maxConcurrent, current)
		}

		// Simulate work
		time.Sleep(50 * time.Millisecond)

		// Decrement concurrent counter
		atomic.AddInt64(&concurrent, -1)

		return item * 2
	})

	result := slices.Collect(mapped.Next)
	elapsed := time.Since(start)

	// Results should be correct
	slices.Sort(result)
	expected := []int{2, 4, 6, 8, 10}
	assert.Equal(t, expected, result)
	assert.NoError(t, mapped.Err())

	// Should have run in parallel (less than sequential time)
	sequentialTime := time.Duration(len(data)) * 50 * time.Millisecond
	assert.Less(t, elapsed, sequentialTime, "Expected parallel execution to be faster than sequential")

	// Should have had multiple concurrent executions
	assert.Greater(t, atomic.LoadInt64(&maxConcurrent), int64(1), "Expected parallel execution")
}

func TestFilterAsync(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	iterator := goiterators.NewIteratorFromSlice(data)

	filtered := goiterators.FilterAsync(iterator, func(item int) bool {
		time.Sleep(10 * time.Millisecond)
		return item%2 == 0
	})

	result := slices.Collect(filtered.Next)
	slices.Sort(result)

	expected := []int{2, 4, 6}
	assert.Equal(t, expected, result)
	assert.NoError(t, filtered.Err())
}

func TestFilterAsyncParallelism(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	iterator := goiterators.NewIteratorFromSlice(data)

	// Track concurrent executions
	var concurrent int64
	var maxConcurrent int64
	var mu sync.Mutex
	processingOrder := make([]int, 0)

	start := time.Now()
	filtered := goiterators.FilterAsync(iterator, func(item int) bool {
		// Track processing order
		mu.Lock()
		processingOrder = append(processingOrder, item)
		mu.Unlock()

		// Increment concurrent counter
		current := atomic.AddInt64(&concurrent, 1)
		if current > atomic.LoadInt64(&maxConcurrent) {
			atomic.StoreInt64(&maxConcurrent, current)
		}

		// Simulate work
		time.Sleep(30 * time.Millisecond)

		// Decrement concurrent counter
		atomic.AddInt64(&concurrent, -1)

		return item%2 == 0
	})

	result := slices.Collect(filtered.Next)
	elapsed := time.Since(start)

	// Results should be correct
	slices.Sort(result)
	expected := []int{2, 4, 6, 8, 10}
	assert.Equal(t, expected, result)
	assert.NoError(t, filtered.Err())

	// Should have run in parallel (less than sequential time)
	sequentialTime := time.Duration(len(data)) * 30 * time.Millisecond
	assert.Less(t, elapsed, sequentialTime, "Expected parallel execution to be faster than sequential")

	// Should have had multiple concurrent executions
	assert.Greater(t, atomic.LoadInt64(&maxConcurrent), int64(1), "Expected parallel execution")

	// Processing order should be different from input order due to parallelism
	mu.Lock()
	assert.NotEqual(t, data, processingOrder, "Expected different processing order due to parallelism")
	mu.Unlock()
}

func TestFilterAsyncEmpty(t *testing.T) {
	data := []int{1, 3, 5}
	iterator := goiterators.NewIteratorFromSlice(data)

	filtered := goiterators.FilterAsync(iterator, func(item int) bool {
		time.Sleep(10 * time.Millisecond)
		return item%2 == 0
	})

	result := slices.Collect(filtered.Next)
	assert.Empty(t, result)
	assert.NoError(t, filtered.Err())
}

func TestMapAsyncWithError(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	next := func(yield func(int, error) bool) {
		for _, item := range data {
			var err error
			if item == 3 {
				err = errors.New("error at 3")
			}
			if !yield(item, err) {
				return
			}
		}
	}

	iterator := goiterators.NewIteratorErr(next)
	mapped := goiterators.MapAsync(iterator, func(item int) int {
		time.Sleep(10 * time.Millisecond)
		return item * 2
	})

	result := slices.Collect(mapped.Next)
	assert.Error(t, mapped.Err())
	// Should only get results before the error
	slices.Sort(result)
	expected := []int{2, 4}
	assert.Equal(t, expected, result)
}

func TestFilterAsyncWithError(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	next := func(yield func(int, error) bool) {
		for _, item := range data {
			var err error
			if item == 4 {
				err = errors.New("error at 4")
			}
			if !yield(item, err) {
				return
			}
		}
	}

	iterator := goiterators.NewIteratorErr(next)
	filtered := goiterators.FilterAsync(iterator, func(item int) bool {
		time.Sleep(10 * time.Millisecond)
		return item%2 == 0
	})

	result := slices.Collect(filtered.Next)
	assert.Error(t, filtered.Err())
	// Should only get results before the error
	slices.Sort(result)
	expected := []int{2}
	assert.Equal(t, expected, result)
}

func TestAsyncErrorPropagationChain(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	next := func(yield func(int, error) bool) {
		for _, item := range data {
			var err error
			if item == 3 {
				err = errors.New("chain error at 3")
			}
			if !yield(item, err) {
				return
			}
		}
	}

	iterator := goiterators.NewIteratorErr(next)

	// Chain async operations
	mapped := goiterators.MapAsync(iterator, func(item int) int {
		time.Sleep(5 * time.Millisecond)
		return item * 2
	})

	filtered := goiterators.FilterAsync(mapped, func(item int) bool {
		time.Sleep(5 * time.Millisecond)
		return item > 2
	})

	result := slices.Collect(filtered.Next)

	// All iterators in the chain should have the same error
	assert.Error(t, iterator.Err())
	assert.Error(t, mapped.Err())
	assert.Error(t, filtered.Err())

	// Should only get results before the error
	slices.Sort(result)
	expected := []int{4} // Only 2*2=4 passes the filter (item > 2)
	assert.Equal(t, expected, result)
}

func TestMapAsyncMidIterationErrorDetection(t *testing.T) {
	// This test tries to create a scenario where an error occurs
	// between iterations, and we want to stop processing immediately

	// Create a slow iterator that sets an error after yielding some items
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	iterator := &slowErrorIterator{
		items:          data,
		errorAfterItem: 3, // Set error after processing item 3
		itemDelay:      20 * time.Millisecond,
	}

	start := time.Now()
	mapped := goiterators.MapAsync(iterator, func(item int) int {
		time.Sleep(10 * time.Millisecond) // Simulate work
		return item * 2
	})

	result := slices.Collect(mapped.Next)
	elapsed := time.Since(start)

	t.Logf("Result: %v", result)
	t.Logf("Error: %v", mapped.Err())
	t.Logf("Elapsed: %v", elapsed)

	// With mid-iteration error checking, we should stop earlier
	// Without it, we process all items and then check for error

	// The key insight: if the iterator is slow and sets an error partway through,
	// we want to detect it quickly rather than waiting for all items

	// Assert that we got the error
	assert.Error(t, mapped.Err())
}

// Custom iterator that sets an error after yielding a specific number of items
type customErrorIterator struct {
	items          []int
	errorAfterItem int
	err            error
	itemsYielded   int
}

func (it *customErrorIterator) Next(yield func(int) bool) {
	for _, item := range it.items {
		if !yield(item) {
			return
		}
		it.itemsYielded++

		// Set error after yielding the specified item
		if it.itemsYielded == it.errorAfterItem {
			it.err = errors.New("error occurred after yielding item")
		}
	}
}

func (it *customErrorIterator) INext(yield func(int, int) bool) {
	for i, item := range it.items {
		if !yield(i, item) {
			return
		}
		it.itemsYielded++

		// Set error after yielding the specified item
		if it.itemsYielded == it.errorAfterItem {
			it.err = errors.New("error occurred after yielding item")
		}
	}
}

func (it *customErrorIterator) Err() error {
	return it.err
}

// Slow iterator that introduces delays between items
type slowErrorIterator struct {
	items          []int
	errorAfterItem int
	itemDelay      time.Duration
	err            error
	itemsYielded   int
}

func (it *slowErrorIterator) Next(yield func(int) bool) {
	for _, item := range it.items {
		if !yield(item) {
			return
		}
		it.itemsYielded++

		// Set error after yielding the specified item
		if it.itemsYielded == it.errorAfterItem {
			it.err = errors.New("error occurred after yielding item")
		}

		// Add delay between items
		time.Sleep(it.itemDelay)
	}
}

func (it *slowErrorIterator) INext(yield func(int, int) bool) {
	for i, item := range it.items {
		if !yield(i, item) {
			return
		}
		it.itemsYielded++

		// Set error after yielding the specified item
		if it.itemsYielded == it.errorAfterItem {
			it.err = errors.New("error occurred after yielding item")
		}

		// Add delay between items
		time.Sleep(it.itemDelay)
	}
}

func (it *slowErrorIterator) Err() error {
	return it.err
}

func TestFlatMapAsync(t *testing.T) {
	data := []int{1, 2, 3}
	iterator := goiterators.NewIteratorFromSlice(data)

	// Each number produces multiple results
	flatMapped := goiterators.FlatMapAsync(iterator, func(item int) iter.Seq[int] {
		time.Sleep(10 * time.Millisecond)
		return slices.Values([]int{item, item * 10})
	})

	result := slices.Collect(flatMapped.Next)
	slices.Sort(result)

	expected := []int{1, 2, 3, 10, 20, 30}
	assert.Equal(t, expected, result)
	assert.NoError(t, flatMapped.Err())
}

func TestForEachAsync(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iterator := goiterators.NewIteratorFromSlice(data)

	var mu sync.Mutex
	var result []int

	// Use ForEachAsync directly
	err := goiterators.ForEachAsync(iterator, func(item int) error {
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		result = append(result, item*2)
		mu.Unlock()
		return nil
	})

	// Check for errors
	assert.NoError(t, err)

	slices.Sort(result)
	expected := []int{2, 4, 6, 8, 10}
	assert.Equal(t, expected, result)
}

func TestForEachAsyncWithError(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iterator := goiterators.NewIteratorFromSlice(data)

	var mu sync.Mutex
	var result []int

	err := goiterators.ForEachAsync(iterator, func(item int) error {
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		result = append(result, item*2)
		mu.Unlock()

		if item == 3 {
			return errors.New("error at item 3")
		}
		return nil
	})

	assert.Error(t, err)
	// Results may vary due to parallelism, but should not be empty
	assert.NotEmpty(t, result)
}

func TestIForEachAsync(t *testing.T) {
	data := []int{10, 20, 30}
	iterator := goiterators.NewIteratorFromSlice(data)

	var mu sync.Mutex
	var result []string

	err := goiterators.IForEachAsync(iterator, func(idx int, item int) error {
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		result = append(result, fmt.Sprintf("idx:%d,val:%d", idx, item))
		mu.Unlock()
		return nil
	})

	slices.Sort(result)
	expected := []string{"idx:0,val:10", "idx:1,val:20", "idx:2,val:30"}
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
}

func TestIForEachAsyncWithError(t *testing.T) {
	data := []int{10, 20, 30}
	iterator := goiterators.NewIteratorFromSlice(data)

	var mu sync.Mutex
	var result []string

	err := goiterators.IForEachAsync(iterator, func(idx int, item int) error {
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		result = append(result, fmt.Sprintf("idx:%d,val:%d", idx, item))
		mu.Unlock()

		if idx == 1 {
			return errors.New("error at index 1")
		}
		return nil
	})

	assert.Error(t, err)
	// Results may vary due to parallelism, but should not be empty
	assert.NotEmpty(t, result)
}

func TestForEachAsyncParallelism(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	iterator := goiterators.NewIteratorFromSlice(data)

	// Track concurrent executions
	var concurrent int64
	var maxConcurrent int64
	var mu sync.Mutex
	var result []int

	start := time.Now()
	err := goiterators.ForEachAsync(iterator, func(item int) error {
		// Increment concurrent counter
		current := atomic.AddInt64(&concurrent, 1)
		if current > atomic.LoadInt64(&maxConcurrent) {
			atomic.StoreInt64(&maxConcurrent, current)
		}

		// Simulate work
		time.Sleep(50 * time.Millisecond)

		// Store result
		mu.Lock()
		result = append(result, item*2)
		mu.Unlock()

		// Decrement concurrent counter
		atomic.AddInt64(&concurrent, -1)

		return nil
	})

	elapsed := time.Since(start)

	assert.NoError(t, err)

	// Results should be correct
	slices.Sort(result)
	expected := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}
	assert.Equal(t, expected, result)

	// Should have run in parallel (less than sequential time)
	sequentialTime := time.Duration(len(data)) * 50 * time.Millisecond
	assert.Less(t, elapsed, sequentialTime, "Expected parallel execution to be faster than sequential")

	// Should have had multiple concurrent executions
	assert.Greater(t, atomic.LoadInt64(&maxConcurrent), int64(1), "Expected parallel execution")
}

func TestForEachAsyncEarlyTermination(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iterator := goiterators.NewIteratorFromSlice(data)

	var mu sync.Mutex
	var result []int

	err := goiterators.ForEachAsync(iterator, func(item int) error {
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		result = append(result, item)
		mu.Unlock()
		
		if item == 3 {
			return errors.New("stopping at 3")
		}
		return nil
	})

	assert.Error(t, err)
	assert.Equal(t, "stopping at 3", err.Error())
	// Due to parallel processing, we might get different numbers of results
	// but we should get at least one result and it should include 3
	assert.NotEmpty(t, result)
	assert.Contains(t, result, 3) // Should contain the item that triggered the error
}

func TestIForEachAsyncEarlyTermination(t *testing.T) {
	data := []int{10, 20, 30, 40, 50}
	iterator := goiterators.NewIteratorFromSlice(data)

	var mu sync.Mutex
	var result []string

	err := goiterators.IForEachAsync(iterator, func(idx int, item int) error {
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		result = append(result, fmt.Sprintf("idx:%d,val:%d", idx, item))
		mu.Unlock()
		
		if idx == 2 {
			return errors.New("stopping at index 2")
		}
		return nil
	})

	assert.Error(t, err)
	assert.Equal(t, "stopping at index 2", err.Error())
	// Due to parallel processing, we might get different numbers of results
	// but we should get at least one result and it should include the error-triggering item
	assert.NotEmpty(t, result)
	
	// Check that the error-triggering item is in the results
	found := false
	for _, res := range result {
		if res == "idx:2,val:30" {
			found = true
			break
		}
	}
	assert.True(t, found, "Should contain the item that triggered the error")
}
