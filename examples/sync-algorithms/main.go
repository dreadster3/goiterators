package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/dreadster3/goiterators"
)

func main() {
	fmt.Println("=== Synchronous Algorithms ===")

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	iter := goiterators.NewIteratorFromSlice(data)

	// Chain operations: filter evens, take first 3, then double each
	result := goiterators.Map(
		goiterators.Take(
			goiterators.Filter(iter, func(x int) bool { return x%2 == 0 }),
			3,
		),
		func(x int) int { return x * 2 },
	)

	output := slices.Collect(result.Next)

	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Filter evens → Take 3 → Double: %v\n", output)

	// String transformation example
	words := []string{"hello", "world", "go"}
	iter2 := goiterators.NewIteratorFromSlice(words)

	uppercase := goiterators.Map(iter2, strings.ToUpper)
	result2 := slices.Collect(uppercase.Next)

	fmt.Printf("\nWords: %v\n", words)
	fmt.Printf("Uppercase: %v\n", result2)
}
