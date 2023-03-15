package src

import (
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TimescaleBenchmark struct {
	logger  *logrus.Entry
	reader  io.Reader
	workers *Workers
}

func NewTimescaleBenchmark(reader io.Reader, size int) *TimescaleBenchmark {
	return &TimescaleBenchmark{
		logger:  logrus.WithTime(time.Now()),
		reader:  reader,
		workers: initWorkers(size),
	}
}

func (tsb *TimescaleBenchmark) Do() {
	tsb.logger.Info("Initializing Timescale Benchmark Assignment")
	tsb.workers.start()
	data := tsb.workers.execute(tsb.createProcesses())
	stats := generateStats(&data)
	fmt.Println(stats.String())
}

func (tsb *TimescaleBenchmark) createProcesses() <-chan Process {

	processes := make(chan Process)

	go func() {

		lines, err := readerToLines(tsb.reader)
		if err != nil {
			tsb.logger.Error(err)
			close(processes)
			return
		}

		for i := 1; i < len(lines); i++ {
			if len(lines[i]) == 0 {
				continue
			}
			row, err2 := lineToRow(string(lines[i]), len(*tsb.workers))
			if err2 != nil {
				tsb.logger.Error(err2)
				continue
			}
			process := Process{
				Id:  uuid.New(),
				Row: row,
			}
			processes <- process
		}

		close(processes)
	}()

	return processes

}
