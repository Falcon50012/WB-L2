package sorter

// Unique удаляет повторяющиеся строки из отсортированного списка
func Unique(records []Record) []Record {
	if len(records) == 0 {
		return records
	}

	j := 0
	for i := 1; i < len(records); i++ {
		if records[i].Line != records[j].Line {
			j++
			records[j] = records[i]
		}
	}

	return records[:j+1]
}
