package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// https://go.dev/blog/pipelines
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	src := stream(ctx)
	worker1 := Sq(ctx, src)
	worker2 := Sq(ctx, src)
	pipeline := Merge(ctx, worker1, worker2)

	for output := range pipeline {
		fmt.Println(output)
	}

	// for i := 0; i < 10; i++ {
	// 	fmt.Println(<-pipeline)
	// }

}

func Producer(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func Mod(input <-chan int, with int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range input {
			out <- n % with
		}
		close(out)
	}()
	return out
}

func Sq(ctx context.Context, input <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range input {
			select {
			case out <- n * n:
			case <-ctx.Done():
				return
			}
		}

	}()
	return out
}

func Sequence(start, end, step int) []int {
	result := make([]int, 0, end-start)
	for i := start; i <= end; i += step {
		result = append(result, i)
	}
	return result
}

func Seq(start, end, step int) <-chan int {
	out := make(chan int)
	go func() {
		for i := start; i <= end; i += step {
			out <- i
		}
		close(out)
	}()
	return out
}

func stream(ctx context.Context) <-chan int {
	out := make(chan int)
	ticker := time.NewTicker(1 * time.Nanosecond)
	cur := 1
	go func() {
		defer close(out)
		for {
			select {
			case <-ticker.C:
				cur++
				out <- cur
			case <-ctx.Done():
				return
			}
		}
	}()
	return out

}

func gen(ctx context.Context, start, end, step int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := start; i <= end; i += step {
			select {
			case out <- i:
				//  a receive operation on a closed channel can always proceed immediately, yielding the element type’s zero value.
			case <-ctx.Done():
				return
			}
		}

	}()
	return out
}

func Merge[T any](ctx context.Context, cs ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan T) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-ctx.Done():
				fmt.Println("Merge canceled")
				return
			}
		}

	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
