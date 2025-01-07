package main

import (
	"fmt"
	"time"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func main() {
	// Goroutines run in the same address space, so access to shared memory *MUST BE SYNCHRONIZED*.  
	// the 'sync' package provides useful primitives, although you won't need them much in GO
	// as there are other primitives better suited for most use cases
	go say("world")
	// go say(..) spins up new goroutine (lightweight thread)
	// EVALUATION of say and params happen in the current go routine (which may be main)
	// EXECUTION of say happens in the new go routine
	fmt.Println("after spun up say(\"world\")")
	say("hello")
	fmt.Println("after say(\"hello\") in main")
}