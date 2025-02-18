package main

import (
	"context"
	"flag"
	"fmt"
	fanin "gcf/basic/boring/patterns/fan-in"
	"gcf/basic/boring/patterns/generator"
	_ "net/http/pprof"
)

var n = flag.Int("n", 10, "max amount of chatter before I get bored")

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	joe := generator.HigherOrderBoring(ctx, "Joe")
	ann := generator.HigherOrderBoring(ctx, "Ann")
	c := fanin.FanIn(ctx, joe, ann)
	for i := 0; i < *n; i++ {
		fmt.Println(<-c)
	}
	fmt.Println("You're both boring; I'm leaving.")
}
