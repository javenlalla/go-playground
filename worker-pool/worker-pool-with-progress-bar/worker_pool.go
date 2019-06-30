package main

import (
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
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
	ProgressBar *mpb.Bar
	sync.Mutex
}

type Task struct {
	Input chan int // Update `Input` to be required channel type.
	Execute func(n int) // Update `Execute` to be required "Work" function signature.
}

func (wp *WorkerPool) Run (done chan bool) {
	var wg sync.WaitGroup
	wg.Add(wp.NumberOfWorkers)

	// Initialize progress bar.
	p := mpb.New(mpb.WithWidth(64))
	name := "Processing:"
	bar := p.AddBar(int64(wp.TotalInput),
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
			decor.OnComplete(
				decor.AverageETA(decor.ET_STYLE_GO, decor.WC{W: 8}), "Complete",
			),
		),
		mpb.AppendDecorators(decor.Percentage()),
	)

	for i := 0; i < wp.NumberOfWorkers; i++ {
		w := Worker{
			Wg: &wg,
			Task: wp.Task,
			TotalInput: wp.TotalInput,
			InputProcesed: &wp.InputProcessed,
			ProgressBar: bar,
		}

		go w.Work()
	}

	wg.Wait()
	p.Wait()
	done <- true
}

func (w *Worker) Work() {
	defer w.Wg.Done()

	for i := range w.Task.Input {
		start := time.Now()
		w.Task.Execute(i)

		w.Lock()
		w.ProgressBar.IncrBy(1, time.Since(start))
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
	for rs := 0; rs < 500; rs++ {
		data = append(data, rs)
	}

	return data
}

// workFunction signature should match the individual element in the `input` channel type.
func workFunction(o int) {
	// Simulate work being done.
	time.Sleep(time.Second)
}
