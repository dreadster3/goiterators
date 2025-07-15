package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/dreadster3/goiterators"
)

func main() {
	fmt.Println("=== Asynchronous Algorithms ===")

	data := []int{1, 2, 3, 4, 5}
	iter := goiterators.NewIteratorFromSlice(data)

	start := time.Now()

	// Async processing - each item processed in parallel
	squared := goiterators.MapAsync(iter, func(x int) int {
		time.Sleep(100 * time.Millisecond) // Simulate expensive work
		fmt.Printf("Processing %d in goroutine\n", x)
		return x * x
	})

	result := slices.Collect(squared.Next)
	elapsed := time.Since(start)

	slices.Sort(result) // Async results may arrive out of order

	fmt.Printf("\nOriginal: %v\n", data)
	fmt.Printf("Squared: %v\n", result)
	fmt.Printf("Time: %v (would be ~500ms if sequential)\n", elapsed)

	// FlatMap example
	iter2 := goiterators.NewIteratorFromSlice([]int{1, 2, 3})
	expanded := goiterators.FlatMapAsync(iter2, func(x int) []int {
		return []int{x, x * 10}
	})

	result2 := slices.Collect(expanded.Next)
	slices.Sort(result2)

	fmt.Printf("\nFlatMap [1,2,3] → [x, x*10]: %v\n", result2)
}
