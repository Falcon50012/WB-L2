package sorter

func Compare(opts Options) func(a, b Record) int {
	return func(a, b Record) int {
		var result int

		switch {
		case opts.Numeric:
			if a.NumOK && b.NumOK {
				result = cmp(a.Num, b.Num)
			} else {
				result = cmp(a.Key, b.Key)
			}

		default:
			result = cmp(a.Key, b.Key)
		}

		if opts.Reverse {
			return -result
		}

		return result
	}
}

func cmp[T string | float64](a, b T) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}
