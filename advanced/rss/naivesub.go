package rss

import "time"

// Convert fetches to a stream
func NaiveSubscribe(fetcher Fetcher) Subscription {
	s := &naivesub{
		fetcher: fetcher,
		updates: make(chan Item), // for Updates
	}
	go s.loop()
	return s
}

// sub implements the Subscription interface.
type naivesub struct {
	fetcher Fetcher   // fetches items
	updates chan Item // delivers items to the user
	err     error
	closed  bool
}

func (s *naivesub) Updates() <-chan Item {
	return s.updates
}

func (s *naivesub) Close() (err error) {
	// TODO: make loop exit
	// TODO: find out about any error
	return err
}

// loop fetches items using s.fetcher and sends them
// on s.updates.  loop exits when s.Close is called.
func (s *naivesub) loop() {
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
