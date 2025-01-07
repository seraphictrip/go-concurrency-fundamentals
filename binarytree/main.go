package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	// We have to close the channel to stop range...
	defer close(ch)
	// use recursive call, but allow for close
	walk(t, ch)

}

// This is an inorder traversal
func walk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	walk(t.Left, ch)
	ch <- t.Value
	walk(t.Right, ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	c1 := make(chan int)
	c2 := make(chan int)
	go Walk(t1, c1)
	go Walk(t2, c2)

	for {
		// NOTE: we need to make sure we close the channels
		v1, ok1 := <-c1
		v2, ok2 := <-c2
		if ok1 != ok2 {
			return false
		}
		if !ok1 {
			// we reached end without having a miss
			return true
		}
		if v1 != v2 {
			return false
		}
	}

}

// There can be many different binary trees with the same sequence of values stored in it. For example,
// here are two binary trees storing the sequence 1, 1, 2, 3, 5, 8, 13.
// see [image](./image/tree.png)
// A function to check whether two binary trees store the same sequence
// is quite comples in most languages (??? - can do inorder traversal...)
// we will use Go's concurrency and channels to write a simple solution.

func main() {
	ch := make(chan int)
	go Walk(tree.New(2), ch)
	for i := range ch {
		fmt.Println(i)
	}
	// true
	fmt.Println(Same(tree.New(1), tree.New(1)))
	// false
	fmt.Println(Same(tree.New(1), tree.New(2)))
	// true
	fmt.Println(Same(tree.New(2), tree.New(2)))
	// false
	fmt.Println(Same(tree.New(2), tree.New(3)))
	// true
	fmt.Println(Same(tree.New(3), tree.New(3)))
	// true
	fmt.Println(Same(tree.New(100), tree.New(100)))
	// false
	fmt.Println(Same(tree.New(10000), tree.New(1000)))
	// true
	fmt.Println(Same(tree.New(10000), tree.New(10000)))
}
