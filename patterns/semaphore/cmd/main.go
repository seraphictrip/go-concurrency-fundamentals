package main

import (
	"fmt"
	"gcf/patterns/semaphore"
	"math/rand"
	"sync"
	"time"
)

func shop(cust int, sem *semaphore.Semaphore, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("cust %d waiting \n", cust)
	sem.Acquire()
	fmt.Printf("cust %d shopping \n", cust)
	time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
	fmt.Printf("cust %d checkinng out \n", cust)
	sem.Release()
}

func main() {
	numCounters := 6
	numCustomers := 1000
	sem := semaphore.NewSemaphore(numCounters)
	wg := new(sync.WaitGroup)
	for i := range numCustomers {
		wg.Add(1)
		go shop(i, sem, wg)
	}

	wg.Wait()
	fmt.Println("shopping frenzy is over")
}
