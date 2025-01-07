package main

import (
	"fmt"
	"sync"
	"time"
)

// Asynchronous worker, the is the process that does our work
// A worker has an ID to identify unit of compute (self)
// An "in" jobs channel (receive only) for the worker to pull jobs from
// and "out" result channel (send only) for the worker to return results
// orr report status
func worker(id int, jobs <-chan int, results chan<- int) {
	// workers will poll jobs from job queue as long as it
	// is no closed, so this task won't end until errors
	// OR closed

	// TEST: panic in a worker
	// if id == 3 {
	// 	// looks like kills whole process, probably needs to be managed
	// 	panic("Does killing a worker kill program?  or just worker?")
	// }

	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		// Process the job
		time.Sleep(250 * time.Millisecond)

		// either return result of work
		// or status if self maintained job
		results <- job * 2
	}
}

func main() {
	// workspace, the work to be done
	numJobs := 1000
	numWorkers := 10

	// In this example we need buffered channels or
	// we might end up in deadlock
	//  need to explore this more and get a better understanding
	// as of now easily gget deadlock of remove buffering of each
	// I'm not sure about if remove from just one
	// I can easily deadlock if set buffer to numWorkers
	// which is what I would know in real world (assuming jobs streaming in)
	// TAKEAWAY: I don't know how to select a good number of buffer
	// I will have to investigate more
	// websearch doesn't really mention deadlock in this scenario
	// maybe some of the other code needs to be explored befoer i jump to conclusions
	//
	jobs := make(chan int, numJobs)
	// the results channel seems to be my bottle neck...
	results := make(chan int, numJobs)

	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			worker(workerID, jobs, results)
		}(i)
	}

	// Enqueue jobs
	// Here we enquue all jobs, this is blocking if channel is full
	// but i have workers already up to pull, so I'm not sure if this is my problem
	for i := 1; i <= numJobs; i++ {
		jobs <- i
	}
	close(jobs)

	func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}
}
