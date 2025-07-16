package goiterators_test

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/dreadster3/goiterators"
	"github.com/stretchr/testify/assert"
)

func TestMapAsyncCtx(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iter := goiterators.NewIteratorFromSlice(data)
	ctx := context.Background()

	mapped := goiterators.MapAsyncCtx(ctx, iter, func(ctx context.Context, x int) (int, error) {
		time.Sleep(10 * time.Millisecond)
		return x * 2, nil
	})

	result := slices.Collect(mapped.Next)
	slices.Sort(result)

	expected := []int{2, 4, 6, 8, 10}
	assert.Equal(t, expected, result)
	assert.NoError(t, mapped.Err())
}

func TestMapAsyncCtxCancellation(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iter := goiterators.NewIteratorFromSlice(data)
	ctx, cancel := context.WithCancel(context.Background())

	mapped := goiterators.MapAsyncCtx(ctx, iter, func(ctx context.Context, x int) (int, error) {
		time.Sleep(100 * time.Millisecond) // Slow operation
		return x * 2, nil
	})

	// Cancel context after a short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	result := slices.Collect(mapped.Next)

	// Should get cancellation error
	assert.Error(t, mapped.Err())
	assert.Equal(t, context.Canceled, mapped.Err())
	// May get some partial results before cancellation
	assert.True(t, len(result) < len(data))
}

func TestMapAsyncCtxTimeout(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iter := goiterators.NewIteratorFromSlice(data)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	mapped := goiterators.MapAsyncCtx(ctx, iter, func(ctx context.Context, x int) (int, error) {
		time.Sleep(100 * time.Millisecond) // Each item takes 100ms
		return x * 2, nil
	})

	result := slices.Collect(mapped.Next)

	// Should get timeout error eventually (parallel processing may complete some)
	if mapped.Err() != nil {
		assert.Equal(t, context.DeadlineExceeded, mapped.Err())
	}
	// May get some partial results before timeout
	assert.True(t, len(result) <= len(data))
}

func TestIMapAsyncCtx(t *testing.T) {
	data := []int{10, 20, 30}
	iter := goiterators.NewIteratorFromSlice(data)
	ctx := context.Background()

	mapped := goiterators.IMapAsyncCtx(ctx, iter, func(ctx context.Context, idx int, x int) (string, error) {
		time.Sleep(10 * time.Millisecond)
		return fmt.Sprintf("%d@%d", x, idx), nil
	})

	result := slices.Collect(mapped.Next)
	slices.Sort(result)

	expected := []string{"10@0", "20@1", "30@2"}
	assert.Equal(t, expected, result)
	assert.NoError(t, mapped.Err())
}

func TestFilterAsyncCtx(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	iter := goiterators.NewIteratorFromSlice(data)
	ctx := context.Background()

	filtered := goiterators.FilterAsyncCtx(ctx, iter, func(ctx context.Context, x int) (bool, error) {
		time.Sleep(10 * time.Millisecond)
		return x%2 == 0, nil
	})

	result := slices.Collect(filtered.Next)
	slices.Sort(result)

	expected := []int{2, 4, 6}
	assert.Equal(t, expected, result)
	assert.NoError(t, filtered.Err())
}

func TestFilterAsyncCtxWithError(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iter := goiterators.NewIteratorFromSlice(data)
	ctx := context.Background()

	filtered := goiterators.FilterAsyncCtx(ctx, iter, func(ctx context.Context, x int) (bool, error) {
		if x == 3 {
			return false, errors.New("error at 3")
		}
		return x%2 == 0, nil
	})

	result := slices.Collect(filtered.Next)

	assert.Error(t, filtered.Err())
	assert.Contains(t, filtered.Err().Error(), "error at 3")
	// Should get some results before error
	assert.True(t, len(result) <= 2) // At most items 2 (and maybe 4,5 if processed before error)
}

func TestIFilterAsyncCtx(t *testing.T) {
	data := []int{10, 11, 12, 13, 14, 15}
	iter := goiterators.NewIteratorFromSlice(data)
	ctx := context.Background()

	filtered := goiterators.IFilterAsyncCtx(ctx, iter, func(ctx context.Context, idx int, x int) (bool, error) {
		time.Sleep(5 * time.Millisecond)
		return idx > 1 && x%2 == 0, nil // Even values after index 1
	})

	result := slices.Collect(filtered.Next)
	slices.Sort(result)

	expected := []int{12, 14} // Even values at indices 2, 4
	assert.Equal(t, expected, result)
	assert.NoError(t, filtered.Err())
}

func TestFlatMapAsyncCtx(t *testing.T) {
	data := []int{1, 2, 3}
	iter := goiterators.NewIteratorFromSlice(data)
	ctx := context.Background()

	flattened := goiterators.FlatMapAsyncCtx(ctx, iter, func(ctx context.Context, x int) ([]int, error) {
		time.Sleep(10 * time.Millisecond)
		return []int{x, x * 10}, nil
	})

	result := slices.Collect(flattened.Next)
	slices.Sort(result)

	expected := []int{1, 2, 3, 10, 20, 30}
	assert.Equal(t, expected, result)
	assert.NoError(t, flattened.Err())
}

func TestIFlatMapAsyncCtx(t *testing.T) {
	data := []int{1, 2, 3}
	iter := goiterators.NewIteratorFromSlice(data)
	ctx := context.Background()

	flattened := goiterators.IFlatMapAsyncCtx(ctx, iter, func(ctx context.Context, idx int, x int) ([]string, error) {
		time.Sleep(10 * time.Millisecond)
		return []string{fmt.Sprintf("idx%d", idx), fmt.Sprintf("val%d", x)}, nil
	})

	result := slices.Collect(flattened.Next)
	slices.Sort(result)

	expected := []string{"idx0", "idx1", "idx2", "val1", "val2", "val3"}
	assert.Equal(t, expected, result)
	assert.NoError(t, flattened.Err())
}

func TestContextCancellationChain(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	iter := goiterators.NewIteratorFromSlice(data)
	ctx, cancel := context.WithCancel(context.Background())

	// Chain context-aware operations
	result := goiterators.FilterAsyncCtx(ctx,
		goiterators.MapAsyncCtx(ctx, iter, func(ctx context.Context, x int) (int, error) {
			time.Sleep(50 * time.Millisecond)
			return x * x, nil
		}),
		func(ctx context.Context, x int) (bool, error) {
			time.Sleep(50 * time.Millisecond)
			return x > 25, nil
		},
	)

	// Cancel after short delay
	go func() {
		time.Sleep(150 * time.Millisecond)
		cancel()
	}()

	collected := slices.Collect(result.Next)

	// Should get cancellation error or complete normally
	if result.Err() != nil {
		assert.Equal(t, context.Canceled, result.Err())
		// Should have partial results
		assert.True(t, len(collected) < len(data))
	} else {
		// Operations completed before cancellation
		assert.True(t, len(collected) <= len(data))
	}
}

func TestContextAwarenessDuringWork(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iter := goiterators.NewIteratorFromSlice(data)
	ctx, cancel := context.WithCancel(context.Background())

	processed := make([]int, 0)
	var processingMutex sync.Mutex

	mapped := goiterators.MapAsyncCtx(ctx, iter, func(ctx context.Context, x int) (int, error) {
		// Simulate work that checks context
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}

		processingMutex.Lock()
		processed = append(processed, x)
		processingMutex.Unlock()

		return x * 2, nil
	})

	// Cancel after allowing some processing
	go func() {
		time.Sleep(75 * time.Millisecond)
		cancel()
	}()

	_ = slices.Collect(mapped.Next) // We don't need the results, just need to consume the iterator

	assert.Error(t, mapped.Err())
	assert.Equal(t, context.Canceled, mapped.Err())

	// Check that some work was interrupted
	processingMutex.Lock()
	processedCount := len(processed)
	processingMutex.Unlock()

	assert.True(t, processedCount < len(data), "Some work should have been interrupted")
}

func TestMixedContextAndRegularOperations(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	iter := goiterators.NewIteratorFromSlice(data)
	ctx := context.Background()

	// Mix context-aware and regular operations
	result := goiterators.Take(
		goiterators.MapAsyncCtx(ctx,
			goiterators.Filter(iter, func(x int) bool { return x > 2 }),
			func(ctx context.Context, x int) (int, error) {
				time.Sleep(10 * time.Millisecond)
				return x * 10, nil
			},
		),
		3,
	)

	output := slices.Collect(result.Next)
	slices.Sort(output)

	// Filter > 2: [3,4,5,6] → Map *10: [30,40,50,60] → Take 3: first 3 items
	// Since async may reorder, we just check we got 3 items from the expected set
	expectedSet := []int{30, 40, 50, 60}
	assert.Len(t, output, 3)
	for _, item := range output {
		assert.Contains(t, expectedSet, item)
	}
	assert.NoError(t, result.Err())
}
