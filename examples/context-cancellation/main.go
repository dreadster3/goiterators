package main

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/dreadster3/goiterators"
)

func main() {
	fmt.Println("=== Context Cancellation Support ===")

	// Example 1: Basic context usage
	fmt.Println("1. Basic context-aware processing:")
	data := []int{1, 2, 3, 4, 5}
	iter := goiterators.NewIteratorFromSlice(data)
	ctx := context.Background()

	result := goiterators.MapAsyncCtx(ctx, iter, func(ctx context.Context, x int) (int, error) {
		time.Sleep(50 * time.Millisecond) // Simulate work
		return x * x, nil
	})

	output := slices.Collect(result.Next)
	slices.Sort(output)
	fmt.Printf("Squared: %v\n", output)

	// Example 2: Context cancellation
	fmt.Println("\n2. Context cancellation:")
	iter2 := goiterators.NewIteratorFromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	ctx2, cancel := context.WithCancel(context.Background())

	start := time.Now()
	result2 := goiterators.MapAsyncCtx(ctx2, iter2, func(ctx context.Context, x int) (int, error) {
		// Simulate long-running work
		time.Sleep(200 * time.Millisecond)
		return x * 10, nil
	})

	// Cancel context after 500ms
	go func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("   🛑 Cancelling context...")
		cancel()
	}()

	output2 := slices.Collect(result2.Next)
	elapsed := time.Since(start)

	fmt.Printf("   Results before cancellation: %v\n", output2)
	fmt.Printf("   Error: %v\n", result2.Err())
	fmt.Printf("   Time elapsed: %v (stopped early)\n", elapsed)

	// Example 3: Context timeout
	fmt.Println("\n3. Context timeout:")
	iter3 := goiterators.NewIteratorFromSlice([]int{1, 2, 3, 4, 5})
	ctx3, cancel3 := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel3()

	start3 := time.Now()
	result3 := goiterators.MapAsyncCtx(ctx3, iter3, func(ctx context.Context, x int) (int, error) {
		time.Sleep(100 * time.Millisecond) // Each item takes 100ms
		return x * 100, nil
	})

	output3 := slices.Collect(result3.Next)
	elapsed3 := time.Since(start3)

	fmt.Printf("   Results before timeout: %v\n", output3)
	fmt.Printf("   Error: %v\n", result3.Err())
	fmt.Printf("   Time elapsed: %v\n", elapsed3)

	// Example 4: Context-aware work with proper cancellation checking
	fmt.Println("\n4. Context-aware work function:")
	iter4 := goiterators.NewIteratorFromSlice([]int{1, 2, 3})
	ctx4, cancel4 := context.WithCancel(context.Background())

	result4 := goiterators.MapAsyncCtx(ctx4, iter4, func(ctx context.Context, x int) (int, error) {
		// Simulate work that can be interrupted
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				return 0, ctx.Err() // Return immediately on cancellation
			default:
				time.Sleep(20 * time.Millisecond) // Do some work
			}
		}
		return x * 1000, nil
	})

	// Cancel after allowing some work
	go func() {
		time.Sleep(150 * time.Millisecond)
		cancel4()
	}()

	output4 := slices.Collect(result4.Next)
	fmt.Printf("   Results: %v\n", output4)
	fmt.Printf("   Error: %v\n", result4.Err())

	// Example 5: Chaining context-aware operations
	fmt.Println("\n5. Chained context operations:")
	iter5 := goiterators.NewIteratorFromSlice([]int{1, 2, 3, 4, 5, 6})
	ctx5, cancel5 := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel5()

	chained := goiterators.FilterAsyncCtx(ctx5,
		goiterators.MapAsyncCtx(ctx5, iter5, func(ctx context.Context, x int) (int, error) {
			time.Sleep(50 * time.Millisecond)
			return x * x, nil
		}),
		func(ctx context.Context, x int) (bool, error) {
			time.Sleep(30 * time.Millisecond)
			return x > 10, nil
		},
	)

	output5 := slices.Collect(chained.Next)
	slices.Sort(output5)
	fmt.Printf("   Chain result: %v\n", output5)
	if err := chained.Err(); err != nil {
		fmt.Printf("   Chain error: %v\n", err)
	}

	// Example 6: Mixed context and regular operations
	fmt.Println("\n6. Mixed context and regular operations:")
	iter6 := goiterators.NewIteratorFromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8})
	ctx6 := context.Background()

	mixed := goiterators.Take(
		goiterators.MapAsyncCtx(ctx6,
			goiterators.Filter(iter6, func(x int) bool { return x%2 == 0 }), // Regular sync filter
			func(ctx context.Context, x int) (int, error) { // Context-aware async map
				time.Sleep(30 * time.Millisecond)
				return x * 10, nil
			},
		),
		3, // Regular sync take
	)

	output6 := slices.Collect(mixed.Next)
	fmt.Printf("   Mixed operations result: %v\n", output6)
	fmt.Printf("   ✓ Seamlessly mix context-aware and regular operations!\n")
}
