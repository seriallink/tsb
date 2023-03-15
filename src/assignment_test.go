package src

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var testParams = `hostname,start_time,end_time
host_000000,2017-01-01 00:33:57,2017-01-01 01:33:57
host_000001,2017-01-01 01:47:29,2017-01-01 02:47:29
host_000002,2017-01-01 01:47:26,2017-01-01 02:47:26
host_000003,2017-01-01 04:30:52,2017-01-01 05:30:52
host_000004,2017-01-01 04:53:35,2017-01-01 05:53:35
host_000005,2017-01-01 00:29:11,2017-01-01 01:29:11
host_000006,2017-01-01 03:18:50,2017-01-01 04:18:50
host_000007,2017-01-01 00:23:38,2017-01-01 01:23:38
host_000008,2017-01-01 05:41:26,2017-01-01 06:41:26
host_000009,2017-01-01 01:45:17,2017-01-01 02:45:17
host_000010,2017-01-01 03:41:42,2017-01-01 04:41:42
host_000011,2017-01-01 06:23:02,2017-01-01 07:23:02
host_000012,2017-01-01 10:27:04,2017-01-01 11:27:04
host_000013,2017-01-01 06:29:05,2017-01-01 07:35:05
host_000014,2017-01-01 07:29:05,2017-01-01 10:29:15`

var testHosts = `hostname,start_time,end_time
host_000000,2017-01-01 00:33:57,2017-01-01 01:33:57
host_000000,2017-01-01 01:27:02,2017-01-01 02:27:02
host_000000,2017-01-01 04:33:34,2017-01-01 05:33:34
host_000001,2017-01-01 01:47:29,2017-01-01 02:47:29
host_000001,2017-01-01 01:56:07,2017-01-01 02:56:07
host_000001,2017-01-01 09:16:10,2017-01-01 10:16:10
host_000002,2017-01-01 01:47:26,2017-01-01 02:47:26
host_000002,2017-01-01 02:17:24,2017-01-01 03:17:24
host_000002,2017-01-01 02:57:58,2017-01-01 03:57:58
host_000003,2017-01-01 04:30:52,2017-01-01 05:30:52
host_000003,2017-01-01 07:53:49,2017-01-01 08:53:49
host_000003,2017-01-01 08:25:22,2017-01-01 09:25:22`

var testLines = `hostname,start_time,end_time
host_000000,2017-01-01 00:33:57,2017-01-01 01:33:57,too many columns
host_000001,too few columns
host_000002,invalid start time,2017-01-01 02:47:26
host_000003,2017-01-01 04:30:52,invalid end time
host_000004,2017-01-01 04:53:35,2017-01-01 05:53:35`

func TestTsbExecution(t *testing.T) {
	err := InitConnectionPool()
	assert.NoError(t, err)
	defer CloseConnectionPool()
	if err == nil {
		tsb := NewTimescaleBenchmark(strings.NewReader(testParams), 5)
		assert.Equal(t, 5, len(*tsb.workers))
		tsb.workers.start()
		data := tsb.workers.execute(tsb.createProcesses())
		data.Range(func(key, value interface{}) bool {
			_, isUUID := key.(uuid.UUID)
			assert.True(t, isUUID)
			_, isResult := value.(Result)
			assert.True(t, isResult)
			return true
		})
		stats := generateStats(&data)
		assert.Equal(t, 15, stats.NumberOfQueries)
		assert.GreaterOrEqual(t, stats.MaximumQueryTime, stats.MinimumQueryTime)
	}
}

func TestWorkerVsHost(t *testing.T) {
	lines, err := readerToLines(strings.NewReader(testHosts))
	assert.NoError(t, err)
	for i := 1; i < len(lines); i++ {
		row, _ := lineToRow(string(lines[i]), 2)
		switch row.Host {
		case "host_000000":
			assert.Equal(t, 0, row.Index)
		case "host_000001":
			assert.Equal(t, 1, row.Index)
		case "host_000002":
			assert.Equal(t, 0, row.Index)
		case "host_000003":
			assert.Equal(t, 1, row.Index)
		}
	}
}

func TestLineToRow(t *testing.T) {
	lines, err := readerToLines(strings.NewReader(testLines))
	assert.NoError(t, err)
	for i := 1; i < len(lines); i++ {
		_, err = lineToRow(string(lines[i]), 1)
		switch i {
		case 1, 2:
			assert.ErrorContains(t, err, "wrong number of columns")
		case 3, 4:
			assert.ErrorContains(t, err, "invalid time format")
		case 5:
			assert.NoError(t, err)
		}
	}
}
