package main

import (
	"log"
	"sync"
	"time"
)

type WorkerPool struct {
	NumberOfWorkers int
	TotalInput int
	InputProcessed int
	Task Task
}

type Worker struct {
	Wg *sync.WaitGroup
	TotalInput int
	InputProcesed *int
	Task Task
	sync.Mutex
}

type Task struct {
	Input chan int // Update `Input` to be required channel type.
	Execute func(n int) // Update `Execute` to be required "Work" function signature.
}

func (wp *WorkerPool) Run (done chan bool) {
	var wg sync.WaitGroup
	wg.Add(wp.NumberOfWorkers)

	for i := 0; i < wp.NumberOfWorkers; i++ {
		w := Worker{
			Wg: &wg,
			Task: wp.Task,
			TotalInput: wp.TotalInput,
			InputProcesed: &wp.InputProcessed,
		}

		go w.Work()
	}

	wg.Wait()
	done <- true
}

func (w *Worker) Work() {
	defer w.Wg.Done()

	for i := range w.Task.Input {
		w.Task.Execute(i)

		w.Lock()

		*w.InputProcesed++
		if *w.InputProcesed % 5 == 0 {
			log.Printf("%d out of %d processed.", *w.InputProcesed, w.TotalInput)
		}

		w.Unlock()
	}
}

func main() {
	input := make(chan int) // Update channel to required channel type.

	data := fetchDataForProcessing()

	done := make(chan bool)
	wp := WorkerPool{
		NumberOfWorkers: 50, // Update number of workers desired.
		TotalInput: len(data),
		Task: Task{
			Input: input,
			Execute: workFunction, // Update target work function.
		},
	}
	go wp.Run(done)

	// Fetch the data needed for processing and send to the Input channel of the WorkPool.
	for _, i := range data {
		input <- i
	}
	close(input)

	<-done
}

func fetchDataForProcessing() []int {
	var data []int
	for rs := 0; rs < 200; rs++ {
		data = append(data, rs)
	}

	return data
}

// workFunction signature should match the individual element in the `input` channel type.
func workFunction(o int) {
	// Simulate work being done.
	time.Sleep(time.Second)
}
