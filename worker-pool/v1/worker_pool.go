package v1

import (
	"fmt"
	"sync"
	"time"
)

type WorkerPool struct {
	NumberOfWorkers int
	Input chan int // Update `Input` to be required channel type.
	WorkFunction func(n int) // Update `Work` to be required function signature.
}

type Worker struct {
	Wg *sync.WaitGroup
	Input chan int // Update `Input` to be required channel type.
	WorkFunction func(n int) // Update `Work` to be required function signature.
}

func (wp *WorkerPool) Run () {
	var wg sync.WaitGroup
	wg.Add(wp.NumberOfWorkers)

	for i := 0; i < wp.NumberOfWorkers; i++ {
		w := Worker{
			Wg: &wg,
			Input: wp.Input,
			WorkFunction: wp.WorkFunction,
		}

		go w.Work()
	}

	wg.Wait()
}

func (w *Worker) Work() {
	defer w.Wg.Done()

	for i := range w.Input {
		w.WorkFunction(i)
	}
}

func main() {
	input := make(chan int) // Update channel to required channel type.

	wp := WorkerPool{
		NumberOfWorkers: 50, // Update number of workers desired.
		Input: input,
		WorkFunction: workFunction, // Update target work function.
	}
	go wp.Run()

	getData(input)
}

// workFunction signature should match the individual element in the `input` channel type.
func workFunction(o int) {
	// Simulate work being done.
	time.Sleep(time.Second)
	fmt.Println(o)
}

func getData(input chan int) {
	for rs := 0; rs < 200; rs++ {
		input <- rs
	}

	close(input)
}
