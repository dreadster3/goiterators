package goiterators_test

import (
	"errors"
	"fmt"
	"iter"
	"slices"
	"testing"

	"github.com/dreadster3/goiterators"
	"github.com/stretchr/testify/assert"
)

func TestTakeSkip(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	next := slices.All(data)

	iterator := goiterators.NewIterator(next)

	mapped := goiterators.Map(iterator, func(item int) int {
		if item > 1 {
			assert.Fail(t, "Not skipping not taken values")
		}

		return item * 2
	})

	take := goiterators.Take(mapped, 1)

	for item := range take.Next {
		t.Log("Item", item)
	}
}

func TestMap(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iterator := goiterators.NewIteratorFromSlice(data)

	mapped := goiterators.Map(iterator, func(item int) int {
		return item * 2
	})

	expected := []int{2, 4, 6, 8, 10}
	result := slices.Collect(mapped.Next)

	assert.Equal(t, expected, result)
	assert.NoError(t, mapped.Err())
}

func TestMapWithError(t *testing.T) {
	data := []int{1, 2, 3}
	next := func(yield func(int, error) bool) {
		for _, item := range data {
			var err error
			if item == 2 {
				err = errors.New("error at 2")
			}
			if !yield(item, err) {
				return
			}
		}
	}

	iterator := goiterators.NewIteratorErr(next)
	mapped := goiterators.Map(iterator, func(item int) int {
		return item * 2
	})

	result := slices.Collect(mapped.Next)
	assert.Error(t, mapped.Err())
	assert.Equal(t, []int{2}, result)
}

func TestFilter(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iterator := goiterators.NewIteratorFromSlice(data)

	filtered := goiterators.Filter(iterator, func(item int) bool {
		return item%2 == 0
	})

	expected := []int{2, 4}
	result := slices.Collect(filtered.Next)

	assert.Equal(t, expected, result)
	assert.NoError(t, filtered.Err())
}

func TestFilterEmpty(t *testing.T) {
	data := []int{1, 3, 5}
	iterator := goiterators.NewIteratorFromSlice(data)

	filtered := goiterators.Filter(iterator, func(item int) bool {
		return item%2 == 0
	})

	result := slices.Collect(filtered.Next)
	assert.Empty(t, result)
	assert.NoError(t, filtered.Err())
}

func TestTake(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iterator := goiterators.NewIteratorFromSlice(data)

	taken := goiterators.Take(iterator, 3)

	expected := []int{1, 2, 3}
	result := slices.Collect(taken.Next)

	assert.Equal(t, expected, result)
	assert.NoError(t, taken.Err())
}

func TestTakeZero(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iterator := goiterators.NewIteratorFromSlice(data)

	taken := goiterators.Take(iterator, 0)

	result := slices.Collect(taken.Next)
	assert.Empty(t, result)
	assert.NoError(t, taken.Err())
}

func TestTakeNegative(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iterator := goiterators.NewIteratorFromSlice(data)

	taken := goiterators.Take(iterator, -1)

	result := slices.Collect(taken.Next)
	assert.Empty(t, result)
	assert.NoError(t, taken.Err())
}

func TestTakeMoreThanAvailable(t *testing.T) {
	data := []int{1, 2, 3}
	iterator := goiterators.NewIteratorFromSlice(data)

	taken := goiterators.Take(iterator, 10)

	result := slices.Collect(taken.Next)
	assert.Equal(t, data, result)
	assert.NoError(t, taken.Err())
}

func TestChainedOperations(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	iterator := goiterators.NewIteratorFromSlice(data)

	filtered := goiterators.Filter(iterator, func(x int) bool {
		return x%2 == 0
	})

	mapped := goiterators.Map(filtered, func(x int) int {
		return x * 2
	})

	taken := goiterators.Take(mapped, 2)

	result := slices.Collect(taken.Next)

	expected := []int{4}
	assert.Equal(t, expected, result)
}

func TestFilterWithError(t *testing.T) {
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
	filtered := goiterators.Filter(iterator, func(item int) bool {
		return item%2 == 0
	})

	result := slices.Collect(filtered.Next)
	assert.Error(t, filtered.Err())
	assert.Equal(t, []int{2}, result)
}

func TestTakeWithError(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	next := func(yield func(int, error) bool) {
		for _, item := range data {
			var err error
			if item == 2 {
				err = errors.New("error at 2")
			}
			if !yield(item, err) {
				return
			}
		}
	}

	iterator := goiterators.NewIteratorErr(next)
	taken := goiterators.Take(iterator, 3)

	result := slices.Collect(taken.Next)
	assert.Error(t, taken.Err())
	assert.Equal(t, []int{1}, result)
}

func TestChainedOperationsWithError(t *testing.T) {
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

	filtered := goiterators.Filter(iterator, func(x int) bool {
		return x%2 == 0
	})

	mapped := goiterators.Map(filtered, func(x int) int {
		return x * 2
	})

	taken := goiterators.Take(mapped, 3)

	result := slices.Collect(taken.Next)

	assert.Error(t, taken.Err())
	assert.Equal(t, []int{4}, result)
}

func TestErrorPropagationOrder(t *testing.T) {
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

	mapped := goiterators.Map(iterator, func(x int) int {
		return x * 2
	})

	taken := goiterators.Take(mapped, 5)

	result := slices.Collect(taken.Next)

	assert.Error(t, iterator.Err())
	assert.Error(t, mapped.Err())
	assert.Error(t, taken.Err())
	assert.Equal(t, []int{2, 4}, result)
}

func TestFlatMap(t *testing.T) {
	data := []int{1, 2, 3}
	iterator := goiterators.NewIteratorFromSlice(data)

	// Each number produces itself and its double
	flattened := goiterators.FlatMap(iterator, func(x int) iter.Seq[int] {
		return slices.Values([]int{x, x * 2})
	})

	result := slices.Collect(flattened.Next)
	expected := []int{1, 2, 2, 4, 3, 6}
	assert.Equal(t, expected, result)
	assert.NoError(t, flattened.Err())
}

func TestFlatMapEmpty(t *testing.T) {
	data := []int{1, 2, 3}
	iterator := goiterators.NewIteratorFromSlice(data)

	// Function returns empty sequence for all items
	flattened := goiterators.FlatMap(iterator, func(x int) iter.Seq[int] {
		return slices.Values([]int{})
	})

	result := slices.Collect(flattened.Next)
	assert.Empty(t, result)
	assert.NoError(t, flattened.Err())
}

func TestFlatMapWithStrings(t *testing.T) {
	data := []string{"ab", "cd"}
	iterator := goiterators.NewIteratorFromSlice(data)

	// Split each string into characters
	flattened := goiterators.FlatMap(iterator, func(s string) iter.Seq[string] {
		chars := make([]string, len(s))
		for i, r := range s {
			chars[i] = string(r)
		}
		return slices.Values(chars)
	})

	result := slices.Collect(flattened.Next)
	expected := []string{"a", "b", "c", "d"}
	assert.Equal(t, expected, result)
	assert.NoError(t, flattened.Err())
}

func TestFlatMapWithError(t *testing.T) {
	// Create iterator with error
	next := func(yield func(int, error) bool) {
		for i := 1; i <= 3; i++ {
			var err error
			if i == 2 {
				err = errors.New("error at 2")
			}
			if !yield(i, err) {
				return
			}
		}
	}

	iterator := goiterators.NewIteratorErr(next)
	flattened := goiterators.FlatMap(iterator, func(x int) iter.Seq[int] {
		return slices.Values([]int{x, x * 10})
	})

	result := slices.Collect(flattened.Next)
	expected := []int{1, 10} // Only first item processed before error
	assert.Equal(t, expected, result)
	assert.Error(t, flattened.Err())
}

func TestForEach(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iterator := goiterators.NewIteratorFromSlice(data)

	var result []int
	err := goiterators.ForEach(iterator, func(item int) error {
		result = append(result, item*2)
		return nil
	})

	expected := []int{2, 4, 6, 8, 10}
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
}

func TestForEachWithError(t *testing.T) {
	next := func(yield func(int, error) bool) {
		for i := 1; i <= 3; i++ {
			var err error
			if i == 2 {
				err = errors.New("error at 2")
			}
			if !yield(i, err) {
				return
			}
		}
	}

	iterator := goiterators.NewIteratorErr(next)
	var result []int
	err := goiterators.ForEach(iterator, func(item int) error {
		result = append(result, item*2)
		return nil
	})

	expected := []int{2} // Only first item processed before error
	assert.Equal(t, expected, result)
	assert.Error(t, err)
}

func TestIForEach(t *testing.T) {
	data := []int{10, 20, 30}
	iterator := goiterators.NewIteratorFromSlice(data)

	var result []string
	err := goiterators.IForEach(iterator, func(idx int, item int) error {
		result = append(result, fmt.Sprintf("idx:%d,val:%d", idx, item))
		return nil
	})

	expected := []string{"idx:0,val:10", "idx:1,val:20", "idx:2,val:30"}
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
}

func TestIForEachWithError(t *testing.T) {
	next := func(yield func(int, error) bool) {
		for i := 1; i <= 3; i++ {
			var err error
			if i == 2 {
				err = errors.New("error at 2")
			}
			if !yield(i, err) {
				return
			}
		}
	}

	iterator := goiterators.NewIteratorErr(next)
	var result []string
	err := goiterators.IForEach(iterator, func(idx int, item int) error {
		result = append(result, fmt.Sprintf("idx:%d,val:%d", idx, item))
		return nil
	})

	expected := []string{"idx:0,val:1"} // Only first item processed before error
	assert.Equal(t, expected, result)
	assert.Error(t, err)
}

func TestForEachEarlyTermination(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iterator := goiterators.NewIteratorFromSlice(data)

	var result []int
	err := goiterators.ForEach(iterator, func(item int) error {
		result = append(result, item)
		if item == 3 {
			return errors.New("stopping at 3")
		}
		return nil
	})

	expected := []int{1, 2, 3} // Should stop after processing 3
	assert.Equal(t, expected, result)
	assert.Error(t, err)
	assert.Equal(t, "stopping at 3", err.Error())
}

func TestIForEachEarlyTermination(t *testing.T) {
	data := []int{10, 20, 30, 40, 50}
	iterator := goiterators.NewIteratorFromSlice(data)

	var result []string
	err := goiterators.IForEach(iterator, func(idx int, item int) error {
		result = append(result, fmt.Sprintf("idx:%d,val:%d", idx, item))
		if idx == 2 {
			return errors.New("stopping at index 2")
		}
		return nil
	})

	expected := []string{"idx:0,val:10", "idx:1,val:20", "idx:2,val:30"} // Should stop after index 2
	assert.Equal(t, expected, result)
	assert.Error(t, err)
	assert.Equal(t, "stopping at index 2", err.Error())
}
