package main

import (
	"context"
	"flag"
	"fmt"
	"gcf/basic/boring/patterns/generator"
)

var n = flag.Int("n", 5, "Number of boring rounds of conversation  before quiting")

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	flag.Parse()
	c := generator.BoringGenerator(ctx, "boring")
	for i := 0; i < 5; i++ {
		fmt.Printf("You say: %q\n", <-c)
	}
	fmt.Println("You're boring; I'm leaving.")

	// Channel as a handle on a "service"
	// our function returns a channel that lets us communicate with the
	// boring service it provides
	// we can have more than one instance
	joe := generator.HigherOrderBoring(ctx, "Joe")
	ann := generator.HigherOrderBoring(ctx, "Ann")
	for i := 0; i < *n; i++ {
		// These are taking turns, not because of actual sequencing, but
		// because of the synchronizing happening here.
		fmt.Println(<-joe)
		fmt.Println(<-ann)
	}

	fmt.Println("You're both boring; I'm leaving.")

}
