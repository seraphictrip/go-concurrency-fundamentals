package main

import (
	"fmt"
	"gcf/basic/boring/patterns/select/daisychain"
)

func transform(secret int) int {
	return secret + 1
}

func main() {
	const n = 1e6
	// we will listen for "secret" on leftmost
	leftmost := make(chan int)
	src := leftmost
	dest := leftmost
	for i := 0; i < n; i++ {
		src = make(chan int)
		go daisychain.Chain(dest, src, transform)
		dest = src
	}
	// start the chain, prior to this no message has come through
	// so everything is blocking waiting for signal
	// code use right <- 1; but start in goroutine to be idomatic?
	go func(c chan int) { c <- 1 }(src)

	fmt.Println(<-leftmost)
}
