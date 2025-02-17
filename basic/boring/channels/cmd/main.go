package main

import (
	"context"
	"flag"
	"fmt"
	"gcf/basic/boring/channels"
	"runtime"
	"time"

	"math/rand"
)

// Our boring examples cheated: the main function couldn't see the output from the other goroutine.

// It was just printed to the screen, where we pretended we saw a conversation.

var n = flag.Int("n", 5, "number of boring messages to receive in main")

// https://go.dev/talks/2012/concurrency.slide#19
// Real conversations require communication.
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	flag.Parse()

	c := channels.MakeUnbufferedChannel[string]()
	go boring(ctx, "boring!", c)
	// recievers should not close channels!
	// so I can't just range over channel...
	// instead I just leave
	for i := 0; i < *n; i++ {
		// Receive expression is just a value (chan is string channel, so a string value)
		// SYNCHRONIZATION: <-c will wait (block) for a value to be sent
		fmt.Printf("You say: %q\n", channels.ReceiveFromChannel(c))
	}
	cancel()
	fmt.Println("You're boring; I'm leaving.")
	time.Sleep(1 * time.Second)
	fmt.Println(runtime.NumGoroutine())
}

func boring(ctx context.Context, msg string, c chan string) {
	defer func() {
		fmt.Println("boring deferred called")
	}()
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			// communicate via the channel
			// send to channel c
			// SYNCHRONIZATION: c <- will wait for a receiver to be ready
			channels.SendToChannel(c, fmt.Sprintf("%s %d", msg, i))
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)

		}

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
