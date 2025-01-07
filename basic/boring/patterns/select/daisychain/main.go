package main

import (
	"fmt"
)

func f(left chan<- int, right <-chan int) {
	secret := <-right
	distortedSecret := secret + 1
	left <- distortedSecret
}

func main() {
	const n = 100000
	// we will listen for "secret" on leftmost
	leftmost := make(chan int)
	right := leftmost
	left := leftmost
	for i := 0; i < n; i++ {
		right = make(chan int)
		go f(left, right)
		left = right
	}
	// start the chain, prior to this no message has come through
	// so everything is blocking waiting for signal
	// code use right <- 1; but start in goroutine to be idomatic?
	go func(c chan int) { c <- 1 }(right)

	fmt.Println(<-leftmost)
}
