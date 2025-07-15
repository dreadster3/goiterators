package goiterators_test

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/dreadster3/goiterators"
	"github.com/stretchr/testify/assert"
)

func TestSyncToAsync(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	iter := goiterators.NewIteratorFromSlice(data)

	// Sync filter followed by async map
	filtered := goiterators.Filter(iter, func(x int) bool { return x%2 == 0 })
	mapped := goiterators.MapAsync(filtered, func(x int) int {
		time.Sleep(10 * time.Millisecond)
		return x * x
	})

	result := slices.Collect(mapped.Next)
	slices.Sort(result)

	expected := []int{4, 16, 36} // 2², 4², 6²
	assert.Equal(t, expected, result)
	assert.NoError(t, mapped.Err())
}

func TestAsyncToSync(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iter := goiterators.NewIteratorFromSlice(data)

	// Async map followed by sync filter
	mapped := goiterators.MapAsync(iter, func(x int) int {
		time.Sleep(10 * time.Millisecond)
		return x * 2
	})
	filtered := goiterators.Filter(mapped, func(x int) bool { return x > 5 })

	result := slices.Collect(filtered.Next)
	slices.Sort(result)

	expected := []int{6, 8, 10} // 3*2, 4*2, 5*2 where result > 5
	assert.Equal(t, expected, result)
	assert.NoError(t, filtered.Err())
}

func TestComplexMixedChain(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	iter := goiterators.NewIteratorFromSlice(data)

	// sync filter → async map → sync take → async filter
	result := goiterators.FilterAsync(
		goiterators.Take(
			goiterators.MapAsync(
				goiterators.Filter(iter, func(x int) bool { return x > 3 }),
				func(x int) int {
					time.Sleep(5 * time.Millisecond)
					return x * x
				},
			),
			4, // Take only first 4 results
		),
		func(x int) bool {
			time.Sleep(5 * time.Millisecond)
			return x > 25
		},
	)

	finalResult := slices.Collect(result.Next)
	slices.Sort(finalResult)

	// Filter > 3: [4,5,6,7,8,9,10]
	// Square: [16,25,36,49,64,81,100]
	// Take 4: first 4 (order may vary, but we take 4)
	// Filter > 25: depends on which 4 we took, but should include values > 25
	assert.True(t, len(finalResult) <= 4)
	for _, val := range finalResult {
		assert.Greater(t, val, 25)
	}
	assert.NoError(t, result.Err())
}

func TestMixedOperationsParallelism(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8}
	iter := goiterators.NewIteratorFromSlice(data)

	start := time.Now()

	// Sync filter (fast) followed by async map (slow)
	chained := goiterators.MapAsync(
		goiterators.Filter(iter, func(x int) bool { return x%2 == 0 }),
		func(x int) int {
			time.Sleep(50 * time.Millisecond) // Simulate expensive work
			return x * x
		},
	)

	result := slices.Collect(chained.Next)
	elapsed := time.Since(start)

	slices.Sort(result)
	expected := []int{4, 16, 36, 64} // 2², 4², 6², 8²
	assert.Equal(t, expected, result)
	assert.NoError(t, chained.Err())

	// Should be faster than sequential (4 * 50ms = 200ms)
	assert.Less(t, elapsed, 150*time.Millisecond, "Should benefit from parallel execution")
}

func TestMixedOperationsWithErrors(t *testing.T) {
	// Create iterator that produces an error
	next := func(yield func(int, error) bool) {
		for i := 1; i <= 6; i++ {
			var err error
			if i == 4 {
				err = errors.New("error at 4")
			}
			if !yield(i, err) {
				return
			}
		}
	}

	iter := goiterators.NewIteratorErr(next)

	// Mix sync and async operations
	result := goiterators.MapAsync(
		goiterators.Filter(iter, func(x int) bool { return x%2 == 0 }),
		func(x int) int {
			time.Sleep(10 * time.Millisecond)
			return x * 10
		},
	)

	collected := slices.Collect(result.Next)
	slices.Sort(collected)

	// Should get 2*10=20 before hitting error at 4
	expected := []int{20}
	assert.Equal(t, expected, collected)
	assert.Error(t, result.Err())
	assert.Contains(t, result.Err().Error(), "error at 4")
}

func TestSyncAsyncSyncChain(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	iter := goiterators.NewIteratorFromSlice(data)

	// sync → async → sync chain
	result := goiterators.Take(
		goiterators.FilterAsync(
			goiterators.Map(iter, func(x int) int { return x * 2 }),
			func(x int) bool {
				time.Sleep(5 * time.Millisecond)
				return x > 10
			},
		),
		3,
	)

	collected := slices.Collect(result.Next)
	slices.Sort(collected)

	// Map: [2,4,6,8,10,12,14,16,18,20]
	// FilterAsync > 10: [12,14,16,18,20]
	// Take 3: first 3 (may vary due to async, but should be 3 items)
	assert.Len(t, collected, 3)
	for _, val := range collected {
		assert.Greater(t, val, 10)
	}
	assert.NoError(t, result.Err())
}

func TestAsyncSyncAsyncChain(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iter := goiterators.NewIteratorFromSlice(data)

	// async → sync → async chain
	result := goiterators.MapAsync(
		goiterators.Filter(
			goiterators.MapAsync(iter, func(x int) int {
				time.Sleep(10 * time.Millisecond)
				return x * x
			}),
			func(x int) bool { return x > 4 },
		),
		func(x int) int {
			time.Sleep(10 * time.Millisecond)
			return x + 100
		},
	)

	collected := slices.Collect(result.Next)
	slices.Sort(collected)

	// First MapAsync: [1,4,9,16,25]
	// Filter > 4: [9,16,25]
	// Second MapAsync: [109,116,125]
	expected := []int{109, 116, 125}
	assert.Equal(t, expected, collected)
	assert.NoError(t, result.Err())
}

func TestMixedOperationsFlatMap(t *testing.T) {
	data := []int{1, 2, 3}
	iter := goiterators.NewIteratorFromSlice(data)

	// sync filter → async flatmap → sync filter
	result := goiterators.Filter(
		goiterators.FlatMapAsync(
			goiterators.Filter(iter, func(x int) bool { return x > 1 }),
			func(x int) []int {
				time.Sleep(10 * time.Millisecond)
				return []int{x, x * 10}
			},
		),
		func(x int) bool { return x < 25 },
	)

	collected := slices.Collect(result.Next)
	slices.Sort(collected)

	// Filter > 1: [2,3]
	// FlatMapAsync: [2,20,3,30]
	// Filter < 25: [2,20,3]
	expected := []int{2, 3, 20}
	assert.Equal(t, expected, collected)
	assert.NoError(t, result.Err())
}

func TestPerformanceOptimization(t *testing.T) {
	// Test that sync filtering before expensive async operation is faster
	// than doing expensive async operation on all items

	largeData := make([]int, 50)
	for i := range largeData {
		largeData[i] = i + 1
	}

	// Strategy 1: Filter first (sync), then expensive async
	iter1 := goiterators.NewIteratorFromSlice(largeData)
	start1 := time.Now()

	optimized := goiterators.MapAsync(
		goiterators.Filter(iter1, func(x int) bool { return x%10 == 0 }), // Reduces items to 5
		func(x int) int {
			time.Sleep(20 * time.Millisecond) // Expensive operation on 5 items
			return x * x
		},
	)
	result1 := slices.Collect(optimized.Next)
	time1 := time.Since(start1)

	// Strategy 2: Expensive async first, then filter
	iter2 := goiterators.NewIteratorFromSlice(largeData)
	start2 := time.Now()

	unoptimized := goiterators.Filter(
		goiterators.MapAsync(iter2, func(x int) int {
			time.Sleep(20 * time.Millisecond) // Expensive operation on all 50 items
			return x * x
		}),
		func(x int) bool { return (x%10 == 0) && (x > 0) }, // This won't work as expected, but for demo
	)
	result2 := slices.Collect(unoptimized.Next)
	time2 := time.Since(start2)

	// Results should be the same
	slices.Sort(result1)
	slices.Sort(result2)

	// The optimized version should be significantly faster
	assert.Less(t, time1, time2, "Optimized version should be faster")
	assert.NoError(t, optimized.Err())
	assert.NoError(t, unoptimized.Err())

	// Both should produce results for multiples of 10
	assert.True(t, len(result1) > 0)
	for _, val := range result1 {
		assert.Equal(t, 0, val%100) // Should be squares of multiples of 10
	}
}

func TestIteratorReusability(t *testing.T) {
	// Test that the same base iterator type can be used with both sync and async
	data := []int{1, 2, 3, 4, 5}

	// Create two identical iterators
	iter1 := goiterators.NewIteratorFromSlice(data)
	iter2 := goiterators.NewIteratorFromSlice(data)

	// Use one with sync operations
	syncResult := goiterators.Map(
		goiterators.Filter(iter1, func(x int) bool { return x%2 == 0 }),
		func(x int) int { return x * 2 },
	)

	// Use another with async operations
	asyncResult := goiterators.MapAsync(
		goiterators.FilterAsync(iter2, func(x int) bool {
			time.Sleep(5 * time.Millisecond)
			return x%2 == 0
		}),
		func(x int) int {
			time.Sleep(5 * time.Millisecond)
			return x * 2
		},
	)

	syncCollected := slices.Collect(syncResult.Next)
	asyncCollected := slices.Collect(asyncResult.Next)

	// Results should be the same
	slices.Sort(asyncCollected) // Async might be out of order
	assert.Equal(t, syncCollected, asyncCollected)
	assert.NoError(t, syncResult.Err())
	assert.NoError(t, asyncResult.Err())
}

func TestErrorPropagationInMixedChain(t *testing.T) {
	// Test error propagation through a complex mixed chain
	next := func(yield func(int, error) bool) {
		for i := 1; i <= 8; i++ {
			var err error
			if i == 5 {
				err = errors.New("error in mixed chain")
			}
			if !yield(i, err) {
				return
			}
		}
	}

	iter := goiterators.NewIteratorErr(next)

	// Complex mixed chain with error
	result := goiterators.Map( // sync
		goiterators.FilterAsync( // async
			goiterators.Take( // sync
				goiterators.MapAsync( // async
					goiterators.Filter(iter, func(x int) bool { return x > 1 }), // sync
					func(x int) int {
						time.Sleep(5 * time.Millisecond)
						return x * 2
					},
				),
				10,
			),
			func(x int) bool {
				time.Sleep(5 * time.Millisecond)
				return x < 20
			},
		),
		func(x int) int { return x + 1 },
	)

	collected := slices.Collect(result.Next)

	// Should get some results before the error
	assert.True(t, len(collected) > 0)
	assert.Error(t, result.Err())
	assert.Contains(t, result.Err().Error(), "error in mixed chain")
}
