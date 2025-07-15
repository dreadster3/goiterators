package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/dreadster3/goiterators"
)

func main() {
	fmt.Println("=== Mixed Sync/Async Operations ===")

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	iter := goiterators.NewIteratorFromSlice(data)

	start := time.Now()

	// Mix sync and async: sync filter → async map → sync take
	result := goiterators.Take(
		goiterators.MapAsync(
			goiterators.Filter(iter, func(x int) bool { return x%2 == 0 }), // Sync: fast
			func(x int) int { // Async: expensive
				time.Sleep(50 * time.Millisecond)
				return x * x
			},
		),
		3, // Sync: limit results
	)

	output := slices.Collect(result.Next)
	elapsed := time.Since(start)

	slices.Sort(output) // Async results may be out of order

	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Filter evens (sync) → Square (async) → Take 3 (sync): %v\n", output)
	fmt.Printf("Time: %v\n", elapsed)
	fmt.Println("\n✓ Same iterator seamlessly uses both sync and async functions!")

	// Performance optimization: cheap operations first
	fmt.Println("\n=== Performance Optimization ===")

	largeData := make([]int, 20)
	for i := range largeData {
		largeData[i] = i + 1
	}

	iter2 := goiterators.NewIteratorFromSlice(largeData)

	start2 := time.Now()
	optimized := goiterators.MapAsync(
		goiterators.Filter(iter2, func(x int) bool { return x%5 == 0 }), // Reduce items first
		func(x int) int {
			time.Sleep(30 * time.Millisecond) // Expensive operation on fewer items
			return x * x
		},
	)

	result2 := slices.Collect(optimized.Next)
	elapsed2 := time.Since(start2)

	slices.Sort(result2)
	fmt.Printf("Filter multiples of 5 → Square: %v\n", result2)
	fmt.Printf("Time: %v (optimized by filtering first)\n", elapsed2)
}
