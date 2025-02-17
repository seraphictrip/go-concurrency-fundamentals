package main

import (
	"fmt"
	"math/rand"
	"time"
)

// https://go.dev/talks/2012/concurrency.slide#12
// This slide just introduces a "boring" example, to avoid any distraction
// this is a synchround  function, that spends most of its time waiting
// introducing the value of concurrency, i.e. we could be doing something
// else during most of this execution time if we had a way to be async
func boring(msg string) {
	for i := 0; ; i++ {
		fmt.Println(msg, i)
		time.Sleep(time.Second)
	}
}

// https://go.dev/talks/2012/concurrency.slide#13
// here we just add some randomness to make it slightly less boring
func slightlyLessBoring(msg string) {
	for i := 0; ; i++ {
		fmt.Println(msg, i)
		// add some randomness to boring convo so doesn't drone on so monotonously
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}

}

// https://go.dev/talks/2012/concurrency.slide#14
// Example 1: not included is just a sync program that runs for ever
/*
func main() {
    boring("boring!")
}
*/

// https://go.dev/talks/2012/concurrency.slide#14
// our first introduction to calling it async, but not very interesting...
// we effectively ignore it because main exits
/*
func main() {
    go boring("boring!")
}
*/

func main() {
	// https://go.dev/talks/2012/concurrency.slide#16
	// A simple example of giving a goroutine time to run..
	// not idiomatic, but start to see the async nature in actions
	// ignore (let happen in own goroutine)
	// great for "Me", but not useful communication
	// i.e. main just exits
	// go slightlyLessBoring("boring!")

	// Ignore a little less
	go slightlyLessBoring("boring!")
	fmt.Println("I'm listening.")
	time.Sleep(2 * time.Second)
	fmt.Println("You're boring; I'm leaving.")
}
