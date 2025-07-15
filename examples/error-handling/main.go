package main

import (
	"errors"
	"fmt"
	"slices"

	"github.com/dreadster3/goiterators"
)

func main() {
	fmt.Println("=== Error Handling ===")

	// Create an iterator that produces an error at item 3
	next := func(yield func(int, error) bool) {
		for i := 1; i <= 5; i++ {
			var err error
			if i == 3 {
				err = errors.New("error at item 3")
			}
			if !yield(i, err) {
				return
			}
		}
	}

	iter := goiterators.NewIteratorErr(next)
	mapped := goiterators.Map(iter, func(x int) int { return x * 2 })

	fmt.Print("Processing: ")
	result := slices.Collect(mapped.Next)

	fmt.Printf("%v\n", result)

	if err := mapped.Err(); err != nil {
		fmt.Printf("Error occurred: %v\n", err)
		fmt.Printf("Got partial results: %v\n", result)
	}

	// Error propagation through async operations
	fmt.Println("\n=== Async Error Propagation ===")

	iter2 := goiterators.NewIteratorErr(next)
	asyncMapped := goiterators.MapAsync(iter2, func(x int) int { return x * 10 })

	result2 := slices.Collect(asyncMapped.Next)
	slices.Sort(result2) // Async results may be out of order

	fmt.Printf("Async results: %v\n", result2)
	if err := asyncMapped.Err(); err != nil {
		fmt.Printf("Async error: %v\n", err)
	}

	// Graceful error handling pattern
	fmt.Println("\n=== Graceful Handling ===")

	processWithGrace := func(iter goiterators.Iterator[int]) {
		results := slices.Collect(iter.Next)

		if err := iter.Err(); err != nil {
			fmt.Printf("⚠️  Processing stopped: %v\n", err)
			fmt.Printf("📊 Partial results: %v\n", results)
		} else {
			fmt.Printf("✅ Success: %v\n", results)
		}
	}

	// With error
	iter3 := goiterators.NewIteratorErr(next)
	processWithGrace(goiterators.Map(iter3, func(x int) int { return x * 5 }))

	// Without error
	cleanIter := goiterators.NewIteratorFromSlice([]int{1, 2, 3})
	processWithGrace(goiterators.Map(cleanIter, func(x int) int { return x * 5 }))
}
