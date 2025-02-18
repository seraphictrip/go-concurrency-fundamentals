package fanin

import "context"

// FanIn combines multiple channels into a single channel
// also may be considered aa multiplexer
/*
// ------\
//		  o--------------
// ------/
*/
func FanIn[T any](ctx context.Context, inputs ...<-chan T) <-chan T {
	mux := make(chan T)
	// for each input channel, spin up a go routine that merges to new stream
	for _, input := range inputs {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					mux <- <-input
				}
			}
		}()
	}
	return mux
}

type Message struct {
	Msg  string
	Wait chan bool
}
