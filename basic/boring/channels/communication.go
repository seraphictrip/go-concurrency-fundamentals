package channels

// A channel in Go provides a connection between two goroutines, allowing them to communicate.

// this is just an unneeded wrapper, it's only purpose is to encapsulate the pattern
// In most cases real code would use impl inline
// INTRODUCE: send-only channel chan<-
func SendToChannel[T any](ch chan<- T, val T) {
	ch <- val
}

// this is just an unneeded wrapper, it's only purpose is to encapsulate the pattern
// In most cases real code would use impl inline
// INTRODUCE: recieve-only channel <-chan
func ReceiveFromChannel[T any](ch <-chan T) T {
	val := <-ch
	return val
}
