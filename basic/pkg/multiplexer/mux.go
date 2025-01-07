package multiplexer

func FanIn[T any](inputs ...<-chan T) <-chan T {
	mux := make(chan T)
	for _, ch := range inputs {
		go func() {
			for {
				val := <-ch
				mux <- val
			}
		}()
	}
	return mux
}

func Muliplex[T any](inputs ...<-chan T) <-chan T {
	return FanIn(inputs...)
}
