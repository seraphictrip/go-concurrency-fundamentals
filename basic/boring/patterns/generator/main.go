package main

import (
	"fmt"
	"time"

	"math/rand"
)

func main() {
	c := boring("boring")
	for i := 0; i < 5; i++ {
		fmt.Printf("You say: %q\n", <-c)
	}
	fmt.Println("You're boring; I'm leaving.")

	// Channel as a handle on a "service"
	// our function returns a channel that lets us communicate with the
	// boring service it provides
	// we can have more than one instance
	joe := boring("Joe")
	ann := boring("Ann")
	for i := 0; i < 5; i++ {
		// These are taking turns, not because of actual sequencing, but
		// because of the synchronizing happening here.
		fmt.Println(<-joe)
		fmt.Println(<-ann)
	}
	fmt.Println("You're both boring; I'm leaving.")

}

// Boring creates it's channel, and returns a receive only channel (of strings)
// Channels are first-class values, just like strings or integers.
func boring(msg string) <-chan string {
	c := make(chan string)
	go func() { // We launch the goroutine from inside the function.
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c // Return the channel to the caller.

}
