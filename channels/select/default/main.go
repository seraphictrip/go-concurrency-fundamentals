package main

import (
	"fmt"
	"time"
)

func main() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			// The default case in a slect is run if no other case is ready
			// USE a default case to try a send or receive wihtout blocking
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// NOTE: use a default case to try a send or receive wihtout blocking
// the above case doesn't not illustrate 
// select {
// case i := <-c:
//     // use i
// default:
//     // receiving from c would block
// }