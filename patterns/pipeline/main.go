package main

import (
	"fmt"
	"time"
)

func main() {
	// Create the initial channel with some data
	// Input
	data := make([]int, 100)
	for i := range data {
		data[i] = i + 1
	}
	// In this example input has to be able to fit all data
	// because it is block in main
	// 1. I should experiment with how to decouple
	input := make(chan int, 10)
	// 1. By simply wrapping in go routine I can shrink buffer size
	go func() {
		for _, d := range data {
			input <- d
		}
		close(input)

	}()
	// First stage of the pipeline: Doubles the input values
	// Stage 1
	doubleOutput := make(chan int)
	go func() {
		defer close(doubleOutput)
		for num := range input {
			doubleOutput <- num * 2
		}
	}()

	// Second stage of the pipeline: Squares the doubled values
	// stage 2
	squareOutput := make(chan int)
	go func() {
		defer close(squareOutput)
		for num := range doubleOutput {
			// pretend this takes time...
			time.Sleep(250 * time.Millisecond)
			squareOutput <- num * num
		}
	}()

	// Third stage of the pipeline: Prints the squared values
	for result := range squareOutput {
		fmt.Println(result)
	}
}
