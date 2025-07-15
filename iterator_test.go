package goiterators_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/dreadster3/goiterators"
	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	iterator := goiterators.NewIteratorFromSlice(data)

	collected := slices.Collect(iterator.Next)

	assert.Equal(t, data, collected)
	assert.NoError(t, iterator.Err())
}

func TestIteratorErr(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	next := func(yield func(int, error) bool) {
		for _, item := range data {
			var err error
			if item == 3 {
				err = errors.New("error")
			}

			if !yield(item, err) {
				return
			}
		}
	}

	iterator := goiterators.NewIteratorErr(next)

	collected := slices.Collect(iterator.Next)
	assert.Error(t, iterator.Err())
	assert.Equal(t, []int{1, 2}, collected)
}

func TestNewIterator(t *testing.T) {
	data := []int{1, 2, 3}
	seq := slices.All(data)

	iterator := goiterators.NewIterator(seq)
	result := slices.Collect(iterator.Next)

	assert.Equal(t, data, result)
	assert.NoError(t, iterator.Err())
}

func TestIteratorINext(t *testing.T) {
	data := []int{10, 20, 30}
	iterator := goiterators.NewIteratorFromSlice(data)

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

func TestIteratorEmpty(t *testing.T) {
	var data []int
	iterator := goiterators.NewIteratorFromSlice(data)

	result := slices.Collect(iterator.Next)
	assert.Empty(t, result)
	assert.NoError(t, iterator.Err())
}

func TestIteratorErrPersistence(t *testing.T) {
	data := []int{1, 2, 3}
	next := func(yield func(int, error) bool) {
		for _, item := range data {
			var err error
			if item == 2 {
				err = errors.New("persistent error")
			}
			if !yield(item, err) {
				return
			}
		}
	}

	iterator := goiterators.NewIteratorErr(next)

	slices.Collect(iterator.Next)
	assert.Error(t, iterator.Err())

	slices.Collect(iterator.Next)
	assert.Error(t, iterator.Err())
}
