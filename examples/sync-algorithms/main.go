package main

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/dreadster3/goiterators"
)

func main() {
	fmt.Println("=== Synchronous Algorithms ===")

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	iterator := goiterators.NewIteratorFromSlice(data)

	// Chain operations: filter evens, take first 3, then double each
	result := goiterators.Map(
		goiterators.Take(
			goiterators.Filter(iterator, func(x int) bool { return x%2 == 0 }),
			3,
		),
		func(x int) int { return x * 2 },
	)

	output := slices.Collect(result.Next)

	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Filter evens → Take 3 → Double: %v\n", output)

	// String transformation example
	words := []string{"hello", "world", "go"}
	iterator2 := goiterators.NewIteratorFromSlice(words)

	uppercase := goiterators.Map(iterator2, strings.ToUpper)
	result2 := slices.Collect(uppercase.Next)

	fmt.Printf("\nWords: %v\n", words)
	fmt.Printf("Uppercase: %v\n", result2)

	// FlatMap example
	fmt.Println("\n=== FlatMap Example ===")
	numbers := []int{1, 2, 3}
	iterator3 := goiterators.NewIteratorFromSlice(numbers)

	// Each number generates itself and its double
	flattened := goiterators.FlatMap(iterator3, func(x int) iter.Seq[int] {
		return slices.Values([]int{x, x * 2})
	})

	result3 := slices.Collect(flattened.Next)
	fmt.Printf("Numbers: %v\n", numbers)
	fmt.Printf("FlatMapped [x, x*2]: %v\n", result3)
}
