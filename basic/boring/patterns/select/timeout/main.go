package main

import (
	"context"
	"fmt"
	"gcf/basic/boring/patterns/generator"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := generator.BoringGenerator(ctx, "joe")
	timeout := time.After(5 * time.Second)
	for {
		select {
		case s := <-c:
			fmt.Println(s)
		// this would re-queue each time, i.e. new timer created on each  loop
		// so if we don't get a message in a second we stop
		case <-time.After(1 * time.Second):
			fmt.Println("You're too slow")
			return
		// this is an overall timeout
		case <-timeout:
			fmt.Println("You talk too much")
			return
		}
	}
}
