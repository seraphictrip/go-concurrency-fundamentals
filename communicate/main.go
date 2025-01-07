// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// @see https://go.dev/doc/codewalk/sharemem/

package main

import (
	"time"
)

const (
	// number of Poller goroutines to launch
	numPollers = 2
	// how often to poll each URL
	pollInterval = 60 * time.Second
	// how often to log status to stdout
	statusInterval = 10 * time.Second
	// back-off timeout on error
	errTimeout = 10 * time.Second
)

// URLs of resources we wish to monitor
var urls = []string{
	"http://www.google.com/",
	"http://golang.org/",
	"http://blog.golang.org/",
	"http://adfadsfa.com/",
}

// Poll Resources on in channel,
func Poller(in <-chan *Resource, out chan<- *Resource, status chan<- State) {
	for r := range in {
		s := r.Poll()
		status <- State{r.url, s, time.Now()}
		out <- r
	}
}

// Don't communicate by sharing memory; share memory by communicating.
func main() {
	/*
		Channels allow you to pass references to data structures between goroutines.
		If you consider this as passing around OWNERSHIP of the data
		(the ability to read and write it),
		they become a powerful and expressive SYNCHRONIZATION mechanism.
	*/
	// Create our input and output channels.
	// Pending = input
	pending, complete := make(chan *Resource), make(chan *Resource)

	// Launch the StateMonitor.
	status := StateMonitor(statusInterval)

	// Launch some Poller goroutines.
	for i := 0; i < numPollers; i++ {
		// Poller will monitor the pending channel
		// extract the resource
		// Poll the resource and update status
		// place resource on completed channel
		// see final loop of main
		go Poller(pending, complete, status)
	}

	// Send some Resources to the pending queue.
	go func() {
		// This is just the inital load of resources
		// we do it in a go routine because pending is an unbuffered channel
		for _, url := range urls {
			pending <- &Resource{url: url}
		}
	}()

	// Each time a Poll is completed it is requeued after
	// a designated interval, adjusted for errors
	for r := range complete {
		// the resource will sleep and then requeue on pending channel
		go r.Sleep(pending)
	}
}
