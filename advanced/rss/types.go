package rss

import (
	"math/rand"
	"time"
)

type Fetcher interface {
	// Fetch from a feed, return items, suggested next fetch time and any error
	Fetch() (items []Item, next time.Time, err error)
}

type FakeFetcher struct {
	domain string
}

func (f *FakeFetcher) Fetch() (items []Item, next time.Time, err error) {
	time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
	return make([]Item, rand.Intn(10)), time.Now().Add(1 * time.Minute), nil
}

// Fetch fetches Items for uri and returns the time when the next
// fetch should be attempted.  On failure, Fetch returns an error.
func Fetch(domain string) Fetcher {
	return &FakeFetcher{domain}
}

type Item struct {
	Title, Channel, GUID string // a subset of RSS fields
}

type Subscription interface {
	Updates() <-chan Item
	Close() error
}

// Convert fetches to a stream
func Subscribe(fetcher Fetcher) Subscription {
	s := &sub{
		fetcher: fetcher,
		updates: make(chan Item), // for Updates
	}
	go s.loop()
	return s
}

// sub implements the Subscription interface.
type sub struct {
	fetcher Fetcher   // fetches items
	updates chan Item // delivers items to the user
	err     error
	closed  bool
}

func (s *sub) Updates() <-chan Item {
	return s.updates
}

func (s *sub) Close() (err error) {
	// TODO: make loop exit
	// TODO: find out about any error
	return err
}

// loop fetches items using s.fetcher and sends them
// on s.updates.  loop exits when s.Close is called.
func (s *sub) loop() {
	for {
		if s.closed {
			close(s.updates)
			return
		}
		items, next, err := s.fetcher.Fetch()
		if err != nil {
			s.err = err
			time.Sleep(10 * time.Second)
			continue
		}
		for _, item := range items {
			s.updates <- item
		}
		if now := time.Now(); next.After(now) {
			time.Sleep(next.Sub(now))
		}
	}
}

// merge serval streams
func Merge(subs ...Subscription) Subscription
