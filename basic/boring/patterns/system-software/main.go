package main

import (
	"fmt"
	"time"

	"math/rand"
)

func main() {
	start := time.Now()
	results := GoogleV3("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}

// The Google function takes a query and returns a slice of Results (which are just strings).
// Google invokes Web, Image, and Video searches serially, appending them to the results slice.
// Run synchronously, so time is total time Web + Image + Video
func GoogleV1(query string) (results []Result) {
	results = append(results, Web(query))
	results = append(results, Image(query))
	results = append(results, Video(query))
	return results
}

// Run the Web, Image, and Video searches concurrently, and wait for all results.
// No locks. No condition variables. No callbacks.
// Run async, so total time is MAX(Web, Image, Video)
func GoogleV2(query string) (results []Result) {
	// Fan-in pattern c (channel) is our multiplexer
	c := make(chan Result)
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	for i := 0; i < 3; i++ {
		result := <-c
		results = append(results, result)
	}
	return results
}

// Like 2.0, but don't wait for slow servers, just send what have as result
// Fan in with timeout
// total is Max(Max(Web, Image, Video), ~80ms), drop unreturned results
func GoogleV2_1(query string) (results []Result) {
	// make channel
	c := make(chan Result)
	// spin up each search
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	// max timeout
	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return results
}

// Reduce tail latency using replicated search servers
func GoogleV3(query string) (results []Result) {
	// fan in
	c := make(chan Result)
	// make lots of requests, but stop on first
	go func() { c <- First(query, Web, Web1, Web2) }()
	go func() { c <- First(query, Image, Image1, Image2) }()
	go func() { c <- First(query, Video, Video1, Video2) }()

	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return results
}

// avoid timeouts by duplicating work..
// spin up multiple calls, and just return the first
func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	searchReplica := func(i int) { c <- replicas[i](query) }
	for i := range replicas {
		go searchReplica(i)
	}

	// as soon as a result comes back, use it
	return <-c
}

var (
	Web    = fakeSearch("web")
	Web1   = fakeSearch("web-1")
	Web2   = fakeSearch("web-2")
	Image  = fakeSearch("image")
	Image1 = fakeSearch("image-1")
	Image2 = fakeSearch("image-2")
	Video  = fakeSearch("video")
	Video1 = fakeSearch("video-1")
	Video2 = fakeSearch("video-2")
)

type Result string

type Search func(query string) Result

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}
