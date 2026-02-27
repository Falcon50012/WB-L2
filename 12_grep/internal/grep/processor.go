package grep

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

func (g *Grep) ProcessFile(w io.Writer, r io.Reader, fileName string, multiple bool) (int, error) {
	reader := bufio.NewReader(r)

	type lineInfo struct {
		num     int
		content []byte
	}

	var beforeBuf []lineInfo
	var blockEnd int
	var lastPrinted int
	var printedInBlock bool

	lineNumber := 0
	matchCount := 0

	if g.cfg.Count {
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil && err != io.EOF {
				return 0, err
			}
			if len(line) == 0 && err == io.EOF {
				break
			}
			lineNumber++
			line = bytes.TrimRight(line, "\r\n")

			matched := g.matcher.Match(line)
			if g.cfg.InvertMatch {
				matched = !matched
			}
			if matched {
				matchCount++
			}

			if err == io.EOF {
				break
			}
		}

		if multiple {
			_, err := fmt.Fprintf(w, "%s:%d\n", fileName, matchCount)
			return matchCount, err
		}
		_, err := fmt.Fprintf(w, "%d\n", matchCount)
		return matchCount, err
	}

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return matchCount, err
		}
		if len(line) == 0 && err == io.EOF {
			break
		}

		lineNumber++
		line = bytes.TrimRight(line, "\r\n")

		matched := g.matcher.Match(line)
		if g.cfg.InvertMatch {
			matched = !matched
		}

		if matched {
			matchCount++

			start := max(lineNumber-g.cfg.Before, 1)

			end := lineNumber + g.cfg.After

			if start <= blockEnd {
				if end > blockEnd {
					blockEnd = end
				}
			} else {
				if printedInBlock {
					fmt.Fprintln(w, "--")
				}
				blockEnd = end
			}

			for _, b := range beforeBuf {
				if b.num >= start && b.num > lastPrinted {
					g.printLine(w, fileName, multiple, b.num, b.content, false)
					lastPrinted = b.num
				}
			}

			if lineNumber > lastPrinted {
				g.printLine(w, fileName, multiple, lineNumber, line, true)
				lastPrinted = lineNumber
			}

			printedInBlock = true
		} else if lineNumber <= blockEnd && lineNumber > lastPrinted {
			g.printLine(w, fileName, multiple, lineNumber, line, false)
			lastPrinted = lineNumber
		}

		if g.cfg.Before > 0 {
			if len(beforeBuf) == g.cfg.Before {
				beforeBuf = beforeBuf[1:]
			}

			beforeBuf = append(beforeBuf, lineInfo{
				num:     lineNumber,
				content: append([]byte{}, line...),
			})
		}

		if err == io.EOF {
			break
		}
	}

	return matchCount, nil
}

func (g *Grep) printLine(w io.Writer, fileName string, multiple bool, num int, line []byte, match bool) {
	if g.cfg.Count {
		return
	}

	if multiple {
		fmt.Fprintf(w, "%s:", fileName)
	}

	if g.cfg.LineNumber {
		if match {
			fmt.Fprintf(w, "%d:%s\n", num, line)
		} else {
			fmt.Fprintf(w, "%d-%s\n", num, line)
		}
	} else {
		fmt.Fprintf(w, "%s\n", line)
	}
}
