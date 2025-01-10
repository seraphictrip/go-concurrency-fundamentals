package main

import (
	"fmt"
	rss "gcf/advanced/rss"
	"time"
)

func main() {
	// Subscribe to some feeds, and create a merged update stream.
	merged := rss.Merge(
		rss.NaiveSubscribe(rss.Fetch("blog.golang.org")),
		rss.NaiveSubscribe(rss.Fetch("googleblog.blogspot.com")),
		rss.NaiveSubscribe(rss.Fetch("googledevelopers.blogspot.com")))

	// Close the subscriptions after some time.
	time.AfterFunc(3*time.Second, func() {
		fmt.Println("closed:", merged.Close())
	})

	// Print the stream.
	for it := range merged.Updates() {
		fmt.Println(it.Channel, it.Title)
	}

	panic("show me the stacks")
}
