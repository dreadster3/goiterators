package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/dreadster3/goiterators"
)

func main() {
	fmt.Println("=== Indexed Operations ===")

	data := []string{"apple", "banana", "cherry", "date", "elderberry"}
	iter := goiterators.NewIteratorFromSlice(data)

	// IMap: Transform with index information
	withIndex := goiterators.IMap(iter, func(idx int, item string) string {
		return fmt.Sprintf("%d:%s", idx, item)
	})

	result1 := slices.Collect(withIndex.Next)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("With index: %v\n", result1)

	// IFilter: Filter based on index
	fmt.Println("\n=== Index-based filtering ===")
	iter2 := goiterators.NewIteratorFromSlice(data)

	evenIndices := goiterators.IFilter(iter2, func(idx int, item string) bool {
		return idx%2 == 0 // Keep items at even indices
	})

	result2 := slices.Collect(evenIndices.Next)
	fmt.Printf("Even indices: %v\n", result2)

	// Regular Take for comparison
	fmt.Println("\n=== Regular Take ===")
	iter3 := goiterators.NewIteratorFromSlice(data)

	first3 := goiterators.Take(iter3, 3)

	result3 := slices.Collect(first3.Next)
	fmt.Printf("First 3: %v\n", result3)

	// Async indexed operations
	fmt.Println("\n=== Async indexed operations ===")
	numbers := []int{1, 2, 3, 4}
	iter4 := goiterators.NewIteratorFromSlice(numbers)

	start := time.Now()
	asyncResult := goiterators.IMapAsync(iter4, func(idx int, item int) string {
		time.Sleep(50 * time.Millisecond) // Simulate work
		return fmt.Sprintf("worker_%d_processed_%d", idx, item*item)
	})

	result4 := slices.Collect(asyncResult.Next)
	elapsed := time.Since(start)

	slices.Sort(result4) // Async results may be out of order
	fmt.Printf("Async results: %v\n", result4)
	fmt.Printf("Time: %v (parallel processing)\n", elapsed)

	// Complex chaining with indexed operations
	fmt.Println("\n=== Chained indexed operations ===")
	iter5 := goiterators.NewIteratorFromSlice([]int{10, 15, 20, 25, 30, 35, 40})

	chained := goiterators.IMap(
		goiterators.IFilter(iter5, func(idx int, item int) bool {
			return idx > 1 && item%5 == 0 // Skip first 2, keep multiples of 5
		}),
		func(idx int, item int) string {
			return fmt.Sprintf("position_%d_value_%d", idx, item)
		},
	)

	result5 := slices.Collect(chained.Next)
	fmt.Printf("Chained result: %v\n", result5)
}
