package main

import "fmt"


// A function defintion, we can call this from any go routine (including main)
func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}

func sum2(s []int) int {
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

func main() {
	s := []int{7, 2, 8, -9, 4, 0}

	// Creating of Channel
	// Like mpas and slices, channels mustbe crated before use
	// channels are typed conduit through which you can send or receive values with the channel operator
	// channel operator: <-
	// ch <- v // Send v to channel ch
	// v := <-ch // receive from ch and assign value to v
	// By default, sends and receives block until the otehr side is ready.  This allows goroutines
	// to synchronize without explicit locks or condition variables

	c := make(chan int)

	// The example code sums the numbers in a slice, 
	// distributing the work between two goroutines. 
	// Once both goroutines have completed their computation, 
	// it calculates the final result.
	// classic dvide and conquer, but concurrently (possible parallel)
	// NOTE: because we are communicating via channels, calls to sum
	// need to be in go routine, else would result in deadlock
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c // receive from c

	fmt.Printf("%v, %v, %v from sum\n", x, y, x+y)

	// sync from sum2
	x = sum2(s[:len(s)/2])
	y = sum2(s[len(s)/2:])
	fmt.Printf("synchronous from sum2 %v, %v, %v\n", x, y, x+y)

	// combine
	go func(){
		c <- sum2(s[:len(s)/2])
	}()
	go func(){
		c <- sum2(s[len(s)/2:])
	}()
	x, y = <-c, <-c // receive from c
	fmt.Printf("from wrapped sum2 %v, %v, %v\n", x, y, x+y)
}