package main

import (
	"fmt"
	"time"

	"math/rand"
)

// Our boring examples cheated: the main function couldn't see the output from the other goroutine.

// It was just printed to the screen, where we pretended we saw a conversation.

// Real conversations require communication.
func main() {
	c := make(chan string)
	go boring("boring!", c)
	// recievers should not close channels!
	// so I can't just range over channel...
	// instead I just leave
	for i := 0; i < 5; i++ {
		// Receive expression is just a value (chan is string channel, so a string value)
		// SYNCHRONIZATION: <-c will wait (block) for a value to be sent
		fmt.Printf("You say: %q\n", <-c)
	}
	fmt.Println("You're boring; I'm leaving.")
}

func boring(msg string, c chan string) {
	for i := 0; ; i++ {
		// communicate via the channel
		// send to channel c
		// SYNCHRONIZATION: c <- will wait for a receiver to be ready
		c <- fmt.Sprintf("%s %d", msg, i)
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

/*
A sender and receiver must both be ready to play their part in the communication. Otherwise we wait until they are.

Thus channels both communicate and synchronize.

ASIDE:
An aside about buffered channels
Note for experts: Go channels can also be created with a buffer.

Buffering removes synchronization.

Buffering makes them more like Erlang's mailboxes.

Buffered channels can be important for some problems but they are more subtle to reason about.

We won't need them today.
*/
