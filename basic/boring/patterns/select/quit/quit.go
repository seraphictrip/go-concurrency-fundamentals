package quit

import (
	"context"
	"fmt"
	"gcf/basic/boring/patterns/generator"
	"time"

	"math/rand"
)

// Flip it around, and have receiver inform sender to stop sending
func BoringWithQuit(ctx context.Context, msg string, quit <-chan bool) <-chan string {
	return generator.RecieveOnlyGenerator(func(ch chan<- string) {
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				return
			case <-quit:
				return
			default:
				ch <- fmt.Sprintf("%s %d", msg, i)
				time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			}
		}
	})
}

func BoringWithBidirectionalQuit(ctx context.Context, msg string, quit chan string) <-chan string {
	return generator.RecieveOnlyGenerator(func(ch chan<- string) {
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				return
			case <-quit:
				// Perform any cleanup...
				// graceful exit
				quit <- "Bye"
				return
			default:
				ch <- fmt.Sprintf("%s %d", msg, i)
				time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			}
		}
	})
}
