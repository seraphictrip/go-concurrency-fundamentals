package main

import "fmt"

func fibonacci(ch, quit chan int) {
	x, y := 0, 1
	for {
		// the *select* statement lets a goroutine wait on multiple communication operations
		// a *select* BLOCKS until one of its cases can run, then it executes that case.
		// It CHOOSES ONE AT RANDOM if multiple are ready.
		select {
		case ch <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	ch := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 32; i++ {
			fmt.Println(<-ch)
			if i == 10 {
				// OK to keep sending on ch
				// we didn't close, but if quit here
				// seems deterministic, always stops at 55
				// quit <- 0
				// if quit in goroutine quit will eventually be called
				// go func() {quit <- 0}()
			}

		}
		// sending to quit here will always first send
		quit <- 0

		// NOTE: If never call quit we would result in deadlock
	}()
	fibonacci(ch, quit)
}