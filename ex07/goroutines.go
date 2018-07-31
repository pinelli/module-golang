package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type Job float64

type Worker int

func (worker Worker) log(msg string) {
	fmt.Printf("worker:%d %s\n", worker, msg)
}

func (worker *Worker) execute(job Job) {
	str := "sleep:" + strconv.FormatFloat(float64(job), 'f', 1, 64)
	worker.log(str)
	time.Sleep(time.Duration(float64(job) * float64(time.Second)))
}

func (worker Worker) stop(workers chan Worker) {
	worker.log("stopping")
	workers <- worker
}

func (worker Worker) run(
	job Job,
	workers chan Worker,
	jobs chan Job) {

	for {
		worker.execute(job)
		more := false
		select {
		case job, more = <-jobs:
			if !more {
				worker.stop(workers)
				return
			}
		default:
			worker.stop(workers)
			return
		}
	}
}

func createWorkers(n int) chan Worker {
	workers := make(chan Worker, n)
	for i := 1; i <= n; i++ {
		workers <- Worker(i)
	}
	return workers
}

func catchAllWorkers(workers chan Worker, n int) {
	for i := 0; i < n-1; i++ {
		<-workers
	}
}

func scheduler(numOfWorkers int, jobs chan Job, wg *sync.WaitGroup) {
	workers := createWorkers(numOfWorkers)
	for {
		worker := <-workers
		job, success := <-jobs

		if !success {
			catchAllWorkers(workers, numOfWorkers)
			wg.Done()
			return
		}

		worker.log("spawning")
		go worker.run(job, workers, jobs)
	}
}

func reader(jobs chan Job) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		txt := scanner.Text()

		var job Job
		_, err := fmt.Sscanf(txt, "%f\n", &job)
		if err == nil {
			jobs <- job
		}
	}
	close(jobs)
}

func stopSignalCheck() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Stdin.Close()
	}()
}

func Run(poolSize int) {
	jobs := make(chan Job, 100)
	var wg sync.WaitGroup
	wg.Add(1) //scheduler

	go stopSignalCheck()
	go scheduler(poolSize, jobs, &wg)
	go reader(jobs)

	wg.Wait()
}

func main() {
	Run(5)
}
