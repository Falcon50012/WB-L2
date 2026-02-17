package sorter

import "strings"

func extractKey(line, delimeter string, column int) string {
	if column <= 0 {
		return line
	}

	if delimeter == "" {
		delimeter = "\t"
	}

	start := 0
	col := 1
	dLen := len(delimeter)

	for col < column {
		idx := strings.Index(line[start:], delimeter)
		if idx == -1 {
			return ""
		}
		start += idx + dLen
		col++
	}

	end := strings.Index(line[start:], delimeter)
	if end == -1 {
		return line[start:]
	}

	return line[start : start+end]
}
