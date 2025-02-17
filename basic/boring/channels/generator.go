package channels

import (
	"fmt"
	"time"

	"math/rand"
)

// Generators are functions that return a channel
// Often one might return the receieve only channel of a src
func BoringGenerator(msg string) <-chan string {
	ch := make(chan string)
	// start up src
	go func() {
		for i := 0; ; i++ {
			ch <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()

	return ch
}

// higher order generator
func RecieveOnlyGenerator[T any](src func(ch chan<- T)) <-chan T {
	ch := make(chan T)
	go src(ch)

	return ch
}
