package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(msg string) {
	for i := 0; ; i++ {
		fmt.Println(msg, i)
		time.Sleep(time.Second)
	}
}

func slightlyLessBoring(msg string) {
	for i := 0; ; i++ {
		fmt.Println(msg, i)
		// add some randomness to boring convo so doesn't drone on so monotonously
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}

}

func main() {
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
