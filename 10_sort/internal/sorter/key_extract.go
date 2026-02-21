package sorter

import "strings"

// extractKey извлекает ключ сортировки из строки
func extractKey(line, delimiter string, column int) string {
	if column <= 0 {
		return line
	}

	if delimiter == "" {
		delimiter = "\t"
	}

	start := 0
	col := 1
	dLen := len(delimiter)

	for col < column {
		idx := strings.Index(line[start:], delimiter)
		if idx == -1 {
			return ""
		}
		start += idx + dLen
		col++
	}

	end := strings.Index(line[start:], delimiter)
	if end == -1 {
		return line[start:]
	}

	return line[start : start+end]
}
