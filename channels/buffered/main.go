package main

import "fmt"

func main() {
	// Channels can be buffered.  Buffer length can be provided as the econd argument
	// to 'make' to initialize a buffered channel.
	ch := make(chan int, 2)
	// doesn't block
	ch <- 1
	// doesn't block
	ch <- 2
	// WOULD BLOCK (possibly deadlock): Sends to a buffered channel block only whewn the buffer is full.
	// ch <- 3
	// doesn't block
	fmt.Println(<-ch)
	// doesn't block
	fmt.Println(<-ch)
	// WOULD BLOCK (possibly deadlock): Receives blcok when the buffer is empty
	// fmt.Println(<-ch)
}