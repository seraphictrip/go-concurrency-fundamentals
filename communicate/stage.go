package main

import (
	"fmt"
	"log"
	"time"
)

// State represents the last-known state of a URL.
type State struct {
	// identifier
	url string
	// status
	status string
	// as of state
	// we poll less often than we report
	// so this is valuable info
	timestamp time.Time
}

// StateMonitor maintains a map that stores the state of the URLs being
// polled, and prints the current state every updateInterval nanoseconds.
// It returns a chan State to which resource state should be sent.
func StateMonitor(updateInterval time.Duration) chan<- State {
	updates := make(chan State)
	urlStatus := make(map[string]string, len(urls))
	// we update on regular interval (10sec by default)
	ticker := time.NewTicker(updateInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				// each update interval log our current status
				logState(urlStatus)
			case s := <-updates:
				// "state" of resources, as string for logging
				urlStatus[s.url] = fmt.Sprintf("%v as of %v", s.status, s.timestamp.Format("15:04:05"))
			}
		}
	}()
	return updates
}

// logState prints a state map.
func logState(s map[string]string) {
	log.Println("Current state:")
	for k, v := range s {
		log.Printf(" %s %s", k, v)
	}
}
