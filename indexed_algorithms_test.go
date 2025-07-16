package goiterators_test

import (
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/dreadster3/goiterators"
	"github.com/stretchr/testify/assert"
)

func TestIMap(t *testing.T) {
	data := []int{10, 20, 30}
	iter := goiterators.NewIteratorFromSlice(data)

	// Use index in transformation
	result := goiterators.IMap(iter, func(idx int, item int) string {
		return fmt.Sprintf("%d@%d", item, idx)
	})

	output := slices.Collect(result.Next)
	expected := []string{"10@0", "20@1", "30@2"}
	assert.Equal(t, expected, output)
	assert.NoError(t, result.Err())
}

func TestIFilter(t *testing.T) {
	data := []int{10, 15, 20, 25, 30}
	iter := goiterators.NewIteratorFromSlice(data)

	// Filter items based on index (keep even indices)
	result := goiterators.IFilter(iter, func(idx int, item int) bool {
		return idx%2 == 0
	})

	output := slices.Collect(result.Next)
	expected := []int{10, 20, 30} // indices 0, 2, 4
	assert.Equal(t, expected, output)
	assert.NoError(t, result.Err())
}


func TestIMapAsync(t *testing.T) {
	data := []int{1, 2, 3}
	iter := goiterators.NewIteratorFromSlice(data)

	result := goiterators.IMapAsync(iter, func(idx int, item int) string {
		time.Sleep(10 * time.Millisecond)
		return fmt.Sprintf("item%d_at_%d", item, idx)
	})

	output := slices.Collect(result.Next)
	slices.Sort(output) // Async results may be out of order

	expected := []string{"item1_at_0", "item2_at_1", "item3_at_2"}
	assert.Equal(t, expected, output)
	assert.NoError(t, result.Err())
}

func TestIFilterAsync(t *testing.T) {
	data := []int{10, 11, 12, 13, 14, 15}
	iter := goiterators.NewIteratorFromSlice(data)

	// Filter async based on both index and value
	result := goiterators.IFilterAsync(iter, func(idx int, item int) bool {
		time.Sleep(5 * time.Millisecond)
		return idx > 1 && item%2 == 0 // Even items after index 1
	})

	output := slices.Collect(result.Next)
	slices.Sort(output) // Async results may be out of order

	expected := []int{12, 14} // Even items at indices 2, 4
	assert.Equal(t, expected, output)
	assert.NoError(t, result.Err())
}

func TestIFlatMapAsync(t *testing.T) {
	data := []int{1, 2, 3}
	iter := goiterators.NewIteratorFromSlice(data)

	result := goiterators.IFlatMapAsync(iter, func(idx int, item int) []string {
		time.Sleep(10 * time.Millisecond)
		return []string{
			fmt.Sprintf("idx%d", idx),
			fmt.Sprintf("val%d", item),
		}
	})

	output := slices.Collect(result.Next)
	slices.Sort(output) // Async results may be out of order

	expected := []string{"idx0", "idx1", "idx2", "val1", "val2", "val3"}
	assert.Equal(t, expected, output)
	assert.NoError(t, result.Err())
}

func TestIndexedChaining(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	iter := goiterators.NewIteratorFromSlice(data)

	// Chain indexed operations
	result := goiterators.IMap(
		goiterators.IFilter(iter, func(idx int, item int) bool {
			return idx%2 == 0 // Even indices
		}),
		func(idx int, item int) string {
			return fmt.Sprintf("pos%d_val%d", idx, item)
		},
	)

	output := slices.Collect(result.Next)
	expected := []string{"pos0_val1", "pos2_val3", "pos4_val5", "pos6_val7", "pos8_val9"}
	assert.Equal(t, expected, output)
	assert.NoError(t, result.Err())
}

func TestMixedIndexedAndRegular(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	iter := goiterators.NewIteratorFromSlice(data)

	// Mix indexed and regular operations
	result := goiterators.Map(
		goiterators.IFilter(iter, func(idx int, item int) bool {
			return idx < 4 // First 4 items
		}),
		func(item int) int { return item * 10 },
	)

	output := slices.Collect(result.Next)
	expected := []int{10, 20, 30, 40}
	assert.Equal(t, expected, output)
	assert.NoError(t, result.Err())
}

func TestIndexedAsyncPerformance(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iter := goiterators.NewIteratorFromSlice(data)

	start := time.Now()
	result := goiterators.IMapAsync(iter, func(idx int, item int) int {
		time.Sleep(50 * time.Millisecond) // Simulate work
		return item*100 + idx
	})

	output := slices.Collect(result.Next)
	elapsed := time.Since(start)

	slices.Sort(output)
	expected := []int{100, 201, 302, 403, 504} // item*100 + idx
	assert.Equal(t, expected, output)
	assert.NoError(t, result.Err())

	// Should be faster than sequential (5 * 50ms = 250ms)
	assert.Less(t, elapsed, 150*time.Millisecond, "Should benefit from parallel execution")
}

func TestIndexedErrorPropagation(t *testing.T) {
	// Create iterator with error
	next := func(yield func(int, error) bool) {
		for i := 1; i <= 5; i++ {
			var err error
			if i == 3 {
				err = errors.New("error at 3")
			}
			if !yield(i, err) {
				return
			}
		}
	}

	iter := goiterators.NewIteratorErr(next)
	result := goiterators.IMap(iter, func(idx int, item int) string {
		return fmt.Sprintf("%d@%d", item, idx)
	})

	output := slices.Collect(result.Next)
	expected := []string{"1@0", "2@1"} // Should stop at error
	assert.Equal(t, expected, output)
	assert.Error(t, result.Err())
}