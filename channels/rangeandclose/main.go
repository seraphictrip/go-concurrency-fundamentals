package main

import (
	"fmt"
)

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	// A SENDER can close a channel to indicate that no more values will be sent.
	// NOTE: Channels are not like files; you don't usually need to close them.
	// Closing is only necessary whwen the receiver must be told there are no more values coming,
	// such as to terminate a 'range' loop.
	close(c)
	// NOTE: Only the sender should close a channel, never the receiver.
	// Sending on a closed channel WILL PANIC
	// c <- 1
}

func main() {
	c := make(chan int, 10)
	go fibonacci(cap(c)*2, c)
	// The loop for i := range receives valeus from the channel repeatedly until it is closed.
	// NOTE: Only the sender should close a channel, never the receiver.
	// Sending on a closed channel WILL PANIC
	for i := range c {
		fmt.Println(i)
	}
	// RECEIVERS can test whether a channel has been closed by assigning a secodn parameter
	// to the receive expression
	v, ok := <-c
	if ok {
		fmt.Printf("received %+v, unexpected....\n", v)
	} else {
		fmt.Println("detected channel closed")
	}
}