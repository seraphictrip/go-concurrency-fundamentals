package main

import (
	"fmt"
	"sync"
	"time"
)

// A Set with a mutex that locks on Add or Has
type SafeSet[T comparable] struct {
	mutex *sync.Mutex
	set   map[T]bool
}

// I know I shouldn't pass mutex as value
// but I am unclear if that means I should be preferring a
// pointer here or not, the tutorial did not
func NewSafeSet[T comparable]() *SafeSet[T] {
	return &SafeSet[T]{
		mutex: &sync.Mutex{},
		set:   make(map[T]bool),
	}
}

func (s *SafeSet[T]) Add(value T) {
	s.mutex.Lock()
	s.set[value] = true
	s.mutex.Unlock()
}

func (s *SafeSet[T]) Has(value T) bool {
	// can run into fatal error: concurrent map read and map write
	// if don't lock
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.set[value]
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
// If I try and pass set by value I get copy lock error
// Crawl passes lock by value: gcf/crawler.SafeSet[string] contains sync.Mutex
// I can fix by passing set by reference, or by setting mutex as reference type in set
// both seem to fix...
// I appreciate the error.
func Crawl(url string, depth int, fetcher Fetcher, set *SafeSet[string]) {
	// Don't fetch the same URL twice.
	if set.Has(url) {
		return
	}
	set.Add(url)

	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	// I'm not sure how to go about this without a Waitgroup.
	// I'm not sure WaitGroups were covered in tutorial, so
	// think maybe I could have done differently???
	wg := &sync.WaitGroup{}
	for _, u := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Crawl(u, depth-1, fetcher, set)
		}()

	}
	wg.Wait()
}

func main() {
	set := NewSafeSet[string]()
	Crawl("https://golang.org/", 4, fetcher, set)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	time.Sleep(time.Duration(1 * time.Millisecond))
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
