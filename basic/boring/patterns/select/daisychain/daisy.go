package daisychain

// Chain channels together
func Chain[T any, V any](dest chan<- T, src <-chan V, transform func(V) T) {
	dest <- transform(<-src)
}
