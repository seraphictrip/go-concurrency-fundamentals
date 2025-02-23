package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu  sync.Mutex
	val int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	c.val++
	c.mu.Unlock()
}

func (c *Counter) Val() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.val
}

// func worker(atom *atomic.Int32, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	atom.Add(1)
// }

func chanWorker(ch chan<- int, delta int, wg *sync.WaitGroup) {
	defer wg.Done()
	ch <- delta
}

func worker(counter *Counter, wg *sync.WaitGroup) {
	defer wg.Done()
	counter.Inc()
}

func main() {
	ch := make(chan int)
	var counter Counter
	var wg sync.WaitGroup
	for range 1000 {
		wg.Add(1)
		go chanWorker(ch, 1, &wg)
		wg.Add(1)
		go worker(&counter, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	acc := 0
	for val := range ch {
		acc += val
	}
	fmt.Println(acc)
	fmt.Println(counter.Val())
}
