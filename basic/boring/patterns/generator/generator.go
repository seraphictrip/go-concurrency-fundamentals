package generator

import (
	"context"
	"fmt"
	"time"

	"math/rand"
)

// Generators are functions that return a channel
// Often one might return the receieve only channel of a src
func BoringGenerator(ctx context.Context, msg string) <-chan string {
	ch := make(chan string)
	// start up src
	go func() {
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- fmt.Sprintf("%s %d", msg, i)
				time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			}
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

// Higher order generator
func HigherOrderBoring(ctx context.Context, msg string) <-chan string {
	return RecieveOnlyGenerator(func(ch chan<- string) {
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				fmt.Println(msg, "canceled")
				return
			default:
				ch <- fmt.Sprintf("%s %d", msg, i)
				time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			}
		}
	})
}
