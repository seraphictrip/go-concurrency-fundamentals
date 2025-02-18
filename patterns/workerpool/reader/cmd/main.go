package main

func main() {

}

func Filter[T any](input <-chan T, predicate func(T) bool) <-chan T {
	out := make(chan T)
	go func() {
		// if input stream closes follow suite
		defer close(out)
		for val := range input {
			if predicate(val) {
				out <- val
			}
		}
	}()

	return out
}

func Merge[T any](buffer int, inputs ...<-chan T) <-chan T {
	out := make(chan T, buffer)
	for _, input := range inputs {
		go func() {
			for val := range input {
				out <- val
			}
		}()
	}
	return out
}
