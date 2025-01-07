package main

import (
	"fmt"
	"time"

	"math/rand"
)

func main() {
	c := FanInGeneric(boring("Joe"), boring("Ann"))
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}
	fmt.Println("You're both boring; I'm leaving.")
}

// fanIn combines multiple channels into a single channel
// also may be considered aa multiplexer
/*
// ------\
//		  o--------------
// ------/
*/
func fanIn(input1, input2 <-chan string) <-chan string {
	// merge/combine
	c := make(chan string)
	go func() {
		for {
			val := <-input1
			c <- val
		}
	}()
	go func() {
		for {
			val := <-input2
			c <- val
		}
	}()
	return c
}

func FanInGeneric[T any](inputs ...<-chan T) <-chan T {
	mux := make(chan T)
	for _, ch := range inputs {
		go func() {
			for {
				val := <-ch
				mux <- val
			}
		}()
	}
	return mux
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
