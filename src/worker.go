package src

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Workers []*Worker

type Worker struct {
	index     int
	logger    *logrus.Entry
	processes chan Process
	results   chan Result
}

type Process struct {
	Id  uuid.UUID
	Row Row
}

type Result struct {
	Id   uuid.UUID
	Span time.Duration
}

func initWorkers(size int) (workers *Workers) {
	workers = new(Workers)
	for i := 0; i < size; i++ {
		*workers = append(*workers, &Worker{
			index:     i,
			processes: make(chan Process),
			results:   make(chan Result),
			logger:    logrus.WithField("worker", i),
		})
	}
	return
}

func (w *Workers) start() {
	for _, worker := range *w {
		go worker.run()
	}
}

func (w *Worker) run() {

	for process := range w.processes {
		w.logger.Debugf("host %s is being processed by worker %d", process.Row.Host, w.index)
		span, err := w.queryTimespan(process.Row)
		if err != nil {
			w.logger.Error(err)
			w.results <- Result{
				Id: process.Id,
			}
			break
		}
		w.results <- Result{
			Id:   process.Id,
			Span: span,
		}
	}

	close(w.results)

}

func (w *Worker) queryTimespan(row Row) (time.Duration, error) {

	sql := `SELECT min(usage), max(usage)                 
              FROM cpu_usage
             WHERE host = $1 
               AND ts BETWEEN $2 AND $3`

	start := time.Now()
	if _, err := GetConnectionPool().Exec(context.Background(), sql, row.Host, row.Start, row.End); err != nil {
		return 0, err
	}

	return time.Now().Sub(start), nil

}

func (w *Workers) execute(processes <-chan Process) (data sync.Map) {

	go func() {
		for process := range processes {
			(*w)[process.Row.Index].processes <- process
		}
		for _, worker := range *w {
			close(worker.processes)
		}
	}()

	results := make(chan Result)

	wg := sync.WaitGroup{}
	wg.Add(len(*w))

	for _, worker := range *w {
		aux := worker // Loop variables captured by func literals in go statements might have unexpected values
		go func() {
			for result := range aux.results {
				results <- result
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		data.Store(r.Id, r)
	}

	return

}
