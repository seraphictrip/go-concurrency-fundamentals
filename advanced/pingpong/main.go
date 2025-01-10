package main

import (
	"fmt"
	"time"
)

type Ball struct{ hits int }

func main() {
	// channels are Typed conduits
	// provide both a way to pass data between go routines
	// and to synchronize ond state
	table := make(chan *Ball)
	go player("ping", table)
	go player("pong", table)

	// if never put ball on table will deadlock
	// can be detected by go runtime
	table <- new(Ball) // game on; toss the ball
	// waiting on main for sleep time, ping pong going on in the background
	time.Sleep(1 * time.Second)
	<-table // game over; grab the ball

	panic("show me the stack")
}

func player(name string, table chan *Ball) {
	for {
		// get ball from table
		ball := <-table
		// increment ball.hit
		ball.hits++
		// inform world what happend
		fmt.Println(name, ball.hits)

		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}
