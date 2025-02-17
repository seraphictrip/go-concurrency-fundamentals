package channels

// Channel Creational patterns

// this is just an unneeded wrapper, it's only purpose is to encapsulate the pattern
// In most cases real code would use impl inline
func MakeUnbufferedChannel[T any]() chan T {
	return make(chan T)
}

// this is just an unneeded wrapper, it's only purpose is to encapsulate the pattern
// In most cases real code would use impl inline
func MakeBufferedChannel[T any](bufsize int) chan T {
	return make(chan T, bufsize)
}
