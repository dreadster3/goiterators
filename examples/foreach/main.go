package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dreadster3/goiterators"
)

func main() {
	fmt.Println("=== ForEach Examples ===")

	// Basic ForEach example
	fmt.Println("\n1. Basic ForEach - Print each element:")
	data := []int{1, 2, 3, 4, 5}
	iter := goiterators.NewIteratorFromSlice(data)

	err := goiterators.ForEach(iter, func(x int) error {
		fmt.Printf("Value: %d\n", x)
		return nil
	})
	if err != nil {
		log.Printf("Error in ForEach: %v", err)
	}

	// IForEach example with index
	fmt.Println("\n2. IForEach - Print with index:")
	fruits := []string{"apple", "banana", "cherry", "date"}
	iter2 := goiterators.NewIteratorFromSlice(fruits)

	err = goiterators.IForEach(iter2, func(idx int, value string) error {
		fmt.Printf("Index %d: %s\n", idx, value)
		return nil
	})
	if err != nil {
		log.Printf("Error in IForEach: %v", err)
	}

	// ForEach with side effects
	fmt.Println("\n3. ForEach with side effects - Calculate sum:")
	numbers := []int{10, 20, 30, 40, 50}
	iter3 := goiterators.NewIteratorFromSlice(numbers)

	sum := 0
	err = goiterators.ForEach(iter3, func(x int) error {
		sum += x
		return nil
	})
	if err != nil {
		log.Printf("Error in ForEach: %v", err)
	} else {
		fmt.Printf("Sum: %d\n", sum)
	}

	// ForEachAsync example
	fmt.Println("\n4. ForEachAsync - Parallel processing:")
	tasks := []int{1, 2, 3, 4, 5}
	iter4 := goiterators.NewIteratorFromSlice(tasks)

	start := time.Now()
	err = goiterators.ForEachAsync(iter4, func(task int) error {
		// Simulate work
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("Processed task %d (goroutine)\n", task)
		return nil
	})
	elapsed := time.Since(start)

	if err != nil {
		log.Printf("Error in ForEachAsync: %v", err)
	} else {
		fmt.Printf("Parallel processing completed in %v\n", elapsed)
	}

	// IForEachAsync example
	fmt.Println("\n5. IForEachAsync - Parallel processing with index:")
	workItems := []string{"task-A", "task-B", "task-C"}
	iter5 := goiterators.NewIteratorFromSlice(workItems)

	start = time.Now()
	err = goiterators.IForEachAsync(iter5, func(idx int, item string) error {
		// Simulate work
		time.Sleep(50 * time.Millisecond)
		fmt.Printf("Completed %s at index %d\n", item, idx)
		return nil
	})
	elapsed = time.Since(start)

	if err != nil {
		log.Printf("Error in IForEachAsync: %v", err)
	} else {
		fmt.Printf("Indexed parallel processing completed in %v\n", elapsed)
	}

	// ForEach with chained operations
	fmt.Println("\n6. ForEach with chained operations:")
	originalData := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	iter6 := goiterators.NewIteratorFromSlice(originalData)

	// Chain operations: filter evens, take 3, then print
	filtered := goiterators.Filter(iter6, func(x int) bool { return x%2 == 0 })
	taken := goiterators.Take(filtered, 3)

	err = goiterators.ForEach(taken, func(x int) error {
		fmt.Printf("Even number: %d\n", x)
		return nil
	})
	if err != nil {
		log.Printf("Error in chained ForEach: %v", err)
	}

	// Early termination example
	fmt.Println("\n7. Early termination - Stop processing when condition is met:")
	terminationData := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	iter7 := goiterators.NewIteratorFromSlice(terminationData)
	
	err = goiterators.ForEach(iter7, func(x int) error {
		fmt.Printf("Processing: %d\n", x)
		if x == 5 {
			return fmt.Errorf("stopping early at %d", x)
		}
		return nil
	})
	
	if err != nil {
		fmt.Printf("Stopped with error: %v\n", err)
	}
	
	// Early termination with async
	fmt.Println("\n8. Early termination (async) - Stop processing when condition is met:")
	terminationData2 := []int{1, 2, 3, 4, 5}
	iter8 := goiterators.NewIteratorFromSlice(terminationData2)
	
	err = goiterators.ForEachAsync(iter8, func(x int) error {
		time.Sleep(20 * time.Millisecond)
		fmt.Printf("Processing (async): %d\n", x)
		if x == 3 {
			return fmt.Errorf("stopping early at %d", x)
		}
		return nil
	})
	
	if err != nil {
		fmt.Printf("Async stopped with error: %v\n", err)
	}

	fmt.Println("\n=== Examples completed successfully! ===")
}
