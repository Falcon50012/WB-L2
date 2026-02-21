package sorter

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

// Read читает строки из файла/Stdin
func Read(filename string, opts Options) ([]Record, error) {
	var (
		f   *os.File
		err error
	)

	if filename != "" {
		f, err = os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
	} else {
		f = os.Stdin
	}

	records := make([]Record, 0, 1024)
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		} else if line == "" && err == io.EOF {
			break
		}

		line = strings.TrimRight(line, "\r\n")

		key := extractKey(line, opts.Delimiter, opts.Column)

		record := Record{
			Line: line,
			Key:  key,
		}

		if opts.Numeric {
			if num, err := strconv.ParseFloat(strings.TrimSpace(key), 64); err == nil {
				record.Num = num
				record.NumOK = true
			}
		}

		records = append(records, record)
	}

	return records, nil
}
