package main

import (
	"context"
	"fmt"
	"gcf/pipeline/producer"
	"slices"
	"time"
)

func main() {
	src := []int{1, 2, 3, 4, 5}

	dest := producer.FromArray(src)
	for val := range dest {
		fmt.Println(val)
	}

	dest = producer.FromIter(slices.Values(src))

	for val := range dest {
		fmt.Println(val)
	}

	count := 0
	ctx, cancel := context.WithCancel(context.Background())
	interval1 := producer.Interval(500*time.Millisecond, ID)
	interval2 := producer.IntervalWithContext(ctx, 100*time.Millisecond, func(_ time.Time) int {
		count++
		return count
	})

	stop := time.After(1 * time.Second)
	stop2 := time.After(10 * time.Second)

	for {
		select {
		case t := <-interval1:
			fmt.Println("interval 1", t)
		case t := <-interval2:
			fmt.Println("interval 2", t)
		case <-stop:
			fmt.Println("cancel")
			cancel()
		case <-stop2:
			return
		default:
			continue
		}
	}

}

func ID(self time.Time) time.Time {
	return self
}
