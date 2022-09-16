package main

import (
	"fmt"
	"sync"
)

type job struct {
	workerIdx int
	jobIdx    int
}

func (j *job) run() {
	fmt.Printf("(worker, job): (%d, %d)\n", j.workerIdx, j.jobIdx)
}

type worker struct {
	workerGroup *sync.WaitGroup
	masterGroup *sync.WaitGroup
}

func (w *worker) work(x *job) {
	w.workerGroup.Add(1)
	defer w.workerGroup.Done()
	x.run()
}

func (w *worker) Dispatch(x *job) {
	go w.work(x)
}

func (w *worker) cleanExit() {
	w.masterGroup.Add(1)
	defer w.masterGroup.Done()
	w.workerGroup.Wait()
}

func (w *worker) Close() error {
	go w.cleanExit()
	return nil
}

type master struct {
	wg *sync.WaitGroup
}

func (m *master) Arrange() *worker {
	return &worker{
		workerGroup: &sync.WaitGroup{},
		masterGroup: m.wg,
	}
}

func (m *master) Close() error {
	m.wg.Wait()
	return nil
}

func main() {
	m := &master{
		wg: &sync.WaitGroup{},
	}
	for i := 0; i < 5; i++ {
		w := m.Arrange()
		for j := 0; j < 10; j++ {
			w.Dispatch(&job{
				workerIdx: i,
				jobIdx:    j,
			})
		}
		w.Close()
	}
	m.Close()
}
