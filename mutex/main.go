package main

import (
	"fmt"
	"sync"
	"time"
)

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

// Inc increments the counter for the given key.
func (c *SafeCounter) Inc(key string) {
	// We can define a block of code to be executed in mutual exclusion
	//  by surrounding it with a call to Lock and Unlock as shown on the Inc method.
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key]++
	c.mu.Unlock()
}

// Value returns the current value of the counter for the given key.
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	// obtain the lock first, before setting up defer
	defer c.mu.Unlock()
	return c.v[key]
}

func main() {
	c := SafeCounter{v: make(map[string]int)}
	for i := 0; i < 1000; i++ {
		go c.Inc("somekey")
		if i%2 == 0 {
			go c.Inc("otherkey")
		}
	}

	time.Sleep(time.Second)
	fmt.Println(c.Value("somekey"))
	fmt.Println(c.Value("otherkey"))
	fmt.Println(c.Value("doesntexist"))
}
