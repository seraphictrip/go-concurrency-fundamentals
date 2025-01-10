package main

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

func main() {
	if len(os.Args) != 2 {
		// usage
		fmt.Fprintln(os.Stderr, "usage: parallel dir")
		return
	}
	m, err := MD5All(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	paths := make([]string, 0, len(m))
	for path := range m {
		paths = append(paths, path)
	}
	slices.Sort(paths)
	for _, path := range paths {
		fmt.Fprintf(os.Stdout, "%x %s\n", m[path], path)
	}
}

// MD5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents.  If the directory walk
// fails or any read operation fails, MD5All returns an error.
func MD5All(root string) (map[string][md5.Size]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// sumFiles is our producer, it fans out and fans in to c
	c, errc := sumFiles(ctx, root)
	m := make(map[string][md5.Size]byte)
	// This is our consumer
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	if err := <-errc; err != nil {
		return nil, err
	}
	return m, nil
}

type result struct {
	// relative path to result
	path string
	sum  [md5.Size]byte // // [16]byte
	// error
	err error
}

func sumFiles(ctx context.Context, root string) (<-chan result, <-chan error) {
	c := make(chan result)
	errc := make(chan error, 1)
	var wg sync.WaitGroup
	go func() {
		// Changed to WalkDir as we can get all the info we need without FileInfo provided by os.Lstat()
		err := filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !info.Type().IsRegular() {
				return nil
			}
			wg.Add(1)
			// spin up (fan out) a go routine per file
			go func() {
				data, err := os.ReadFile(path)
				select {
				// fan-in to c
				case c <- result{path, md5.Sum(data), err}:
				case <-ctx.Done():
				}
				wg.Done()
			}()
			// abort the wlak if done is closed
			select {
			case <-ctx.Done():
				return errors.New("walk cancelled")
			default:
				return nil
			}
		})
		// Walk has returned, so all calls to wg.Add are done.  Start a
		// goroutine to close c once all the sends are done.
		go func() {
			wg.Wait()
			close(c)
		}()
		// No select needed here, since errc is buffered.
		// we only put one error (possibly nil) on this chan
		errc <- err
	}()

	return c, errc
}
