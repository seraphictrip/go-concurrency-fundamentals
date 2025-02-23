package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i + 1
	}
	// as with other example, we putting this in buffer first
	// so need enough room in buffer
	// we can break the dependency by simply wrapping input into a go routine
	input := make(chan int, len(data))

	go func() {
		for _, d := range data {
			input <- d
		}
		close(input)
	}()

	// Fan-out: Launch multiple worker goroutines
	numWorkers := 3
	results := make(chan int, numWorkers)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		// add to waitgroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for num := range input {
				// Simulate some processing
				time.Sleep(10 * time.Millisecond)
				result := num * 2
				// Fan-in: Aggregate results from workers
				results <- result
			}
		}()
	}

	go func() {
		// if this is not in a go routine then I would wait to start processing
		wg.Wait()
		close(results)
	}()

	// Process aggregated results
	for result := range results {
		fmt.Println(result)
	}
}
