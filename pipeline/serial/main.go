package main

// See https://go.dev/blog/pipeline serial.go
import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

// No concurrency, this implementation simply reads and sums each file as it walks the tree
func main() {
	if len(os.Args) != 2 {
		// usage
		fmt.Fprintln(os.Stderr, "usage: serial dir")
		return
	}
	// MD5All does all the heavy lifting
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
	m := make(map[string][md5.Size]byte)
	// Changed to WalkDir as we can get all the info we need without FileInfo provided by os.Lstat()
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.Type().IsRegular() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		m[path] = md5.Sum(data)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}
