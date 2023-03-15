package src

import (
	"bytes"
	"fmt"
	"sort"
	"sync"
	"time"
)

type SummaryStats struct {
	NumberOfQueries     int
	TotalProcessingTime time.Duration
	MinimumQueryTime    time.Duration
	MaximumQueryTime    time.Duration
	AverageQueryTime    time.Duration
	MedianQueryTime     time.Duration
}

func (s *SummaryStats) String() string {

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("number of queries processed: %d\n", s.NumberOfQueries))
	buffer.WriteString(fmt.Sprintf("total processing time across all queries: %v\n", s.TotalProcessingTime))
	buffer.WriteString(fmt.Sprintf("minimum query time (for a single query): %v\n", s.MinimumQueryTime))
	buffer.WriteString(fmt.Sprintf("maximum query time: %v\n", s.MaximumQueryTime))
	buffer.WriteString(fmt.Sprintf("median query time: %v\n", s.MedianQueryTime))
	buffer.WriteString(fmt.Sprintf("average query time: %v\n", s.AverageQueryTime))

	return buffer.String()

}

func calculateMedian(durations []time.Duration) time.Duration {

	length := len(durations)
	mid := length / 2

	if length%2 == 0 {
		return (durations[mid-1] + durations[mid]) / 2
	}

	return durations[mid]

}

func generateStats(results *sync.Map) (stats SummaryStats) {

	durations := make([]time.Duration, 0, 0)

	results.Range(func(key, value interface{}) bool {
		stats.NumberOfQueries += 1
		stats.TotalProcessingTime += value.(Result).Span
		durations = append(durations, value.(Result).Span)
		return true
	})

	if len(durations) == 0 {
		return stats
	}

	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	stats.MinimumQueryTime = durations[0]
	stats.MaximumQueryTime = durations[len(durations)-1]
	stats.AverageQueryTime = time.Duration(int(stats.TotalProcessingTime) / stats.NumberOfQueries)
	stats.MedianQueryTime = calculateMedian(durations)

	return stats

}
