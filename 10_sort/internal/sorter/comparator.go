package sorter

func Compare(opts Options) func(a, b Record) int {
	return func(a, b Record) int {
		var result int

		switch {
		case opts.Numeric:
			if a.NumOK && b.NumOK {
				result = numeric(a.Num, b.Num)
			} else {
				result = lexical(a.Key, b.Key)
			}

		default:
			result = lexical(a.Key, b.Key)
		}

		if opts.Reverse {
			return -result
		}

		return result
	}
}

func lexical(a, b string) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func numeric(a, b float64) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}
