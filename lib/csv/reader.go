package csv

import (
	"bufio"
	csv "encoding/csv"
	"os"
)

func newCsvReader(file os.File) *csv.Reader {
	reader := csv.NewReader(file)
	br := bufio.NewReader(r)
	bs, err := br.Peek(3)
	if err != nil {
		return csv.NewReader(br)
	}
	if bs[0] == 0xEF && bs[1] == 0xBB && bs[2] == 0xBF {
		br.Discard(3)
	}
	return csv.NewReader(br)
}
