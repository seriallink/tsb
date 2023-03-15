package src

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"io"
	"strings"
	"time"
)

const dateTimeLayout = "2006-01-02 15:04:05"

type Row struct {
	Index int
	Host  string
	Start time.Time
	End   time.Time
}

func getIndexByHost(host string, size int) int {
	h := fnv.New64a()
	h.Write([]byte(host))
	return int(h.Sum64() % uint64(size))
}

func readerToLines(reader io.Reader) ([][]byte, error) {

	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return bytes.Split(b, []byte{'\n'}), nil

}

func lineToRow(line string, size int) (row Row, err error) {

	values := strings.Split(line, ",")
	if len(values) != 3 {
		return Row{}, fmt.Errorf("wrong number of columns: %v", line)
	}

	row.Host = values[0]
	row.Index = getIndexByHost(row.Host, size)

	if row.Start, err = time.Parse(dateTimeLayout, values[1]); err != nil {
		return Row{}, fmt.Errorf("invalid time format: %s", values[1])
	}

	if row.End, err = time.Parse(dateTimeLayout, values[2]); err != nil {
		return Row{}, fmt.Errorf("invalid time format: %s", values[2])
	}

	return

}
