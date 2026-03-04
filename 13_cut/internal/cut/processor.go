package cut

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
)

// isFieldSelected — бинарный поиск по нормализованным диапазонам.
func isFieldSelected(ranges []rangeSpec, idx int) bool {
	if len(ranges) == 0 {
		return false
	}
	i := sort.Search(len(ranges), func(i int) bool {
		return ranges[i].start > idx
	})
	if i == 0 {
		return false
	}
	r := ranges[i-1]
	if r.end == -1 {
		return idx >= r.start
	}
	return idx >= r.start && idx <= r.end
}

// joinWithDelim склеивает поля через delim.
func joinWithDelim(fields [][]byte, delim []byte) []byte {
	if len(fields) == 0 {
		return []byte{}
	}
	if len(fields) == 1 {
		return fields[0]
	}

	total := 0
	for _, f := range fields {
		total += len(f)
	}
	total += (len(fields) - 1) * len(delim)

	b := make([]byte, 0, total)
	for i, f := range fields {
		if i > 0 {
			b = append(b, delim...)
		}
		b = append(b, f...)
	}
	return b
}

// ProcessFile читает из r построчно, разбивает по разделителю и выводит выбранные поля в w.
func (c *Cut) ProcessFile(w io.Writer, r io.Reader) (err error) {
	if len(c.delim) == 0 {
		return errors.New("пустой разделитель")
	}

	br := bufio.NewReader(r)
	bw := bufio.NewWriter(w)

	defer func() {
		if e := bw.Flush(); err == nil {
			err = e
		} else if e != nil {
			fmt.Fprintln(os.Stderr, "flush error:", e)
		}
	}()

	oneByteDelim := len(c.delim) == 1
	var delimByte byte
	if oneByteDelim {
		delimByte = c.delim[0]
	}
	delimLen := len(c.delim)

	for {
		lineBytes, readErr := br.ReadBytes('\n')
		if readErr != nil && readErr != io.EOF {
			return readErr
		}

		hasNewline := len(lineBytes) > 0 && lineBytes[len(lineBytes)-1] == '\n'

		trimmed := bytes.TrimRight(lineBytes, "\r\n")

		var firstPos int
		if oneByteDelim {
			firstPos = bytes.IndexByte(trimmed, delimByte)
		} else {
			firstPos = bytes.Index(trimmed, c.delim)
		}

		if firstPos == -1 {
			if c.cfg.OnlyDelimited {
				if readErr == io.EOF {
					break
				}
				continue
			}

			if isFieldSelected(c.ranges, 1) {
				if _, werr := bw.Write(trimmed); werr != nil {
					return werr
				}
			}

			if hasNewline {
				if werr := bw.WriteByte('\n'); werr != nil {
					return werr
				}
			}

			if readErr == io.EOF {
				break
			}

			continue
		}

		var outFields [][]byte
		idx := 1

		field := trimmed[:firstPos]
		if isFieldSelected(c.ranges, idx) {
			outFields = append(outFields, field)
		}

		start := firstPos + delimLen
		idx++

		for start <= len(trimmed) {
			var pos int

			if oneByteDelim {
				pos = bytes.IndexByte(trimmed[start:], delimByte)
			} else {
				pos = bytes.Index(trimmed[start:], c.delim)
			}

			if pos == -1 {
				field := trimmed[start:]

				if isFieldSelected(c.ranges, idx) {
					outFields = append(outFields, field)
				}

				break
			}

			field := trimmed[start : start+pos]

			if isFieldSelected(c.ranges, idx) {
				outFields = append(outFields, field)
			}

			start += pos + delimLen
			idx++
		}

		if len(outFields) > 0 {
			if _, werr := bw.Write(joinWithDelim(outFields, c.delim)); werr != nil {
				return werr
			}
		}

		if hasNewline {
			if werr := bw.WriteByte('\n'); werr != nil {
				return werr
			}
		}

		if readErr == io.EOF {
			break
		}
	}

	return nil
}
