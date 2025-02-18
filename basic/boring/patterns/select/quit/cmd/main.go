package main

import (
	"context"
	"fmt"
	"gcf/basic/boring/patterns/select/quit"

	"math/rand"
)

// flip direction, and we can have a way to tell a channel to quit
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	quitChan := make(chan bool)
	c := quit.BoringWithQuit(ctx, "joe", quitChan)
	for i := rand.Intn(10); i >= 0; i-- {
		fmt.Println(<-c)
	}
	fmt.Println("stop talking!")
	quitChan <- true
	// we might want to wait for them to respond,
	// if for example sender needs to do cleanup for graceful exit
	biQuit := make(chan string)
	c = quit.BoringWithBidirectionalQuit(ctx, "steve", biQuit)
	for i := rand.Intn(10); i >= 0; i-- {
		fmt.Println(<-c)
	}
	fmt.Println("stop talking!")
	biQuit <- "shut up"

	fmt.Println("steve says", <-biQuit)
}
