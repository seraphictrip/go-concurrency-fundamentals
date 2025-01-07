package main

import (
	"fmt"
	"gcf/basic/pkg/multiplexer"

	"time"

	"math/rand"
)

func main() {
	waitForIt := make(chan bool)
	c := multiplexer.FanIn(boring("Joe", waitForIt), boring("Ann", waitForIt))
	for i := 0; i < 5; i++ {
		msg1 := <-c
		fmt.Println(msg1.str)
		msg2 := <-c
		fmt.Println(msg2.str)
		msg1.wait <- true
		msg2.wait <- true
	}
}

// Channels are first class citizens, you can pass channels to chanels
type Message struct {
	str string
	// signaler channel
	//
	wait chan bool
}

func boring(msg string, waitForIt chan bool) <-chan Message {
	c := make(chan Message)
	go func() { // We launch the goroutine from inside the function.
		for i := 0; ; i++ {
			c <- Message{fmt.Sprintf("%s: %d", msg, i), waitForIt}
			time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
			<-waitForIt
		}
	}()
	return c // Return the channel to the caller.
}
