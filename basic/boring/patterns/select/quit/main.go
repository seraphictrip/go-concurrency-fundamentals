package main

import (
	"fmt"

	"math/rand"
)

// flip direction, and we can have a way to tell a channel to quit
func main() {
	quit := make(chan string)
	c := boring("Joe", quit)
	for i := rand.Intn(10); i >= 0; i-- {
		fmt.Println(<-c)
	}
	fmt.Println("stop talking!")
	quit <- "bye"
	// we might want to wait for them to respond,
	// if for example they need to do cleanup
	fmt.Printf("Joe says: %q\n", <-quit)
}

func boring(msg string, quit chan string) <-chan string {
	c := make(chan string)
	go func() { // We launch the goroutine from inside the function.
		for i := 0; ; i++ {
			select {
			case c <- fmt.Sprintf("%s: %d", msg, i):
				// do nothing
			case <-quit:
				// do any cleanup()
				quit <- "See ya!"
				return
			}
		}
	}()
	return c // Return the channel to the caller.
}
