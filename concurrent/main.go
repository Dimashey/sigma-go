package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Task struct {
	// Name represents the task's name
	Name string
	// Duration represents the time the task should take to execute
	Duration time.Duration
	// Result will store the computed result of the task after completion
	Result int
}

func RunTask(t *Task, r chan Task) {
	fmt.Printf("%s started\n", t.Name)
	time.Sleep(t.Duration)

	t.Result = int(t.Duration.Seconds() * 1000)

	r <- *t
}

func FanOutTasks(tasks []Task, wg *sync.WaitGroup, results chan<- Task, ctx context.Context) {
	for _, task := range tasks {
		wg.Add(1)

		go func(t *Task) {
			defer wg.Done()
			r := make(chan Task, 1)

			go RunTask(t, r)

			for {
				select {
				case <-ctx.Done():
					fmt.Printf("%s canceled due to timeout\n", t.Name)
					close(results)
					return
				case result := <-r:
					results <- result
					return
				default:
					time.Sleep(time.Second)
					fmt.Printf("Working on task %s\n", t.Name)
				}
			}
		}(&task)
	}
	wg.Wait()
}

func FanInResults(results <-chan Task) {
	for t := range results {
		fmt.Printf("%s completed with result: %d\n", t.Name, t.Result)
	}
}

func TimeoutManager(tasks []Task, timeout time.Duration) {
	results := make(chan Task, len(tasks))
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup

	FanOutTasks(tasks, &wg, results, ctx)
	FanInResults(results)
}

func main() {
	tasks := []Task{
		{Name: "Task 1", Duration: 5 * time.Second},
		{Name: "Task 2", Duration: 2 * time.Second},
		{Name: "Task 3", Duration: 1 * time.Second},
		{Name: "Task 4", Duration: 1 * time.Second},
		{Name: "Task 5", Duration: 1 * time.Second},
		{Name: "Task 6", Duration: 1 * time.Second},
	}

	TimeoutManager(tasks, 3*time.Second)
}
