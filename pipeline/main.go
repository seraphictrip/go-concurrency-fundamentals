package main

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io/fs"
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

	c, errc := sumFiles(ctx, root)
	m := make(map[string][md5.Size]byte)
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
		err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			wg.Add(1)
			go func() {
				data, err := os.ReadFile(path)
				select {
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
		errc <- err
	}()

	return c, errc
}
