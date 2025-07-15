package main

import (
	"fmt"
	"slices"

	"github.com/dreadster3/goiterators"
)

func main() {
	fmt.Println("=== Basic Iterator Usage ===")

	// Create an iterator from a slice
	numbers := []int{1, 2, 3, 4, 5}
	iter := goiterators.NewIteratorFromSlice(numbers)

	fmt.Printf("Original: %v\n", numbers)

	// Iterate and collect results
	result := slices.Collect(iter.Next)
	fmt.Printf("Collected: %v\n", result)

	// Using indices
	iter2 := goiterators.NewIteratorFromSlice([]string{"hello", "world"})
	fmt.Println("\nWith indices:")
	for idx, item := range iter2.INext {
		fmt.Printf("%d: %s\n", idx, item)
	}
}
