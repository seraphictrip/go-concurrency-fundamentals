package main

import (
	"fmt"
	"time"

	"math/rand"
)

func main() {

	c := boring("Joe")
	timeout := time.After(5 * time.Second)
	for {
		select {
		case s := <-c:
			fmt.Println(s)
		// this would re-queue each time
		// so if we don't get a message in a second we stop
		case <-time.After(1 * time.Second):
			fmt.Println("You're too slow")
			return
		// this is an overall timeout
		case <-timeout:
			fmt.Println("You talk too much")
			return
		}
	}
}

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
