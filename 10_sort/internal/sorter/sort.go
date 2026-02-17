package sorter

import "slices"

func Sort(records []Record, opts Options) {
	slices.SortStableFunc(records, Compare(opts))
}
