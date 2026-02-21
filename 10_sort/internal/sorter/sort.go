package sorter

import "slices"

// Sort выполняет стабильную сортировку строк
func Sort(records []Record, opts Options) {
	slices.SortStableFunc(records, Compare(opts))
}
