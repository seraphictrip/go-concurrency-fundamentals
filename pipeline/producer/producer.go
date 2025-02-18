package producer

import (
	"context"
	"iter"
	"time"
)

// Produce from an Array
func FromArray[T any](src []T) <-chan T {
	output := make(chan T)

	go func() {
		defer close(output)
		for _, val := range src {
			output <- val
		}
	}()
	return output
}

// Produce from an Iterable
func FromIter[T any](src iter.Seq[T]) <-chan T {
	output := make(chan T)

	go func() {
		defer close(output)
		for val := range src {
			output <- val
		}
	}()

	return output
}

func Interval[T any](d time.Duration, mapper func(time.Time) T) <-chan T {
	ticker := time.NewTicker(d)
	ch := make(chan T)

	go func() {
		defer close(ch)
		for val := range ticker.C {
			ch <- mapper(val)
		}
	}()

	return ch
}

func IntervalWithContext[T any](ctx context.Context, d time.Duration, mapper func(time.Time) T) <-chan T {
	ticker := time.NewTicker(d)
	ch := make(chan T)

	go func() {
		defer close(ch)
		for {
			select {
			case t := <-ticker.C:
				ch <- mapper(t)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	return ch
}

func CountInterval(ctx context.Context, interval time.Duration) <-chan int {
	count := 0
	return IntervalWithContext(ctx, interval, func(_ time.Time) int {
		count++
		return count
	})
}
