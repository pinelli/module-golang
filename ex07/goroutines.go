package goroutines

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Job float64

type Worker struct {
	id int
}

func (worker *Worker) log(msg string) {
	fmt.Printf("worker:%d %s\n", worker.id, msg)
}

func (worker *Worker) execute(job Job) {
	str := "sleep:" + strconv.FormatFloat(float64(job), 'f', 1, 64)
	worker.log(str)
	time.Sleep(time.Duration(float64(job)*1000) * time.Millisecond)
}

func (worker *Worker) run(job Job, workers chan *Worker, jobs chan Job) {
	for {
		worker.execute(job)
		select {
		case job = <-jobs:
		default:
			worker.log("stopping")
			workers <- worker
			return
		}
	}
}

func createWorkers(n int) chan *Worker {
	workers := make(chan *Worker, n)
	for i := 1; i <= n; i++ {
		workers <- &Worker{i}
	}
	return workers
}

func scheduler(numOfWorkers int, jobs chan Job) {
	workers := createWorkers(numOfWorkers)
	for {
		worker := <-workers
		job := <-jobs
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
		fmt.Sscanf(txt, "%f\n", &job)

		jobs <- job
	}
}

func Run(poolSize int) {
	var jobs = make(chan Job, 100)

	go scheduler(poolSize, jobs)
	go reader(jobs)
}
