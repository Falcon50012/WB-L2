package cut

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// rangeSpec описывает интервал полей [start..end].
type rangeSpec struct {
	start int
	end   int // -1 => infinity
}

// Cut — основной контекст утилиты.
type Cut struct {
	cfg    Config
	ranges []rangeSpec
	delim  []byte
}

// NewCut парсит конфигурацию и подготавливает нормализованные диапазоны.
func NewCut(cfg Config) (*Cut, error) {
	if strings.TrimSpace(cfg.Fields) == "" {
		return nil, errors.New("поля должны быть указаны (флаг -f)")
	}

	parsed, err := parseFieldsSpec(cfg.Fields)
	if err != nil {
		return nil, fmt.Errorf("недопустимая конфигурация полей: %w", err)
	}

	if cfg.Delimiter == "" {
		cfg.Delimiter = "\t"
	}

	if strings.ContainsAny(cfg.Delimiter, "\r\n") {
		return nil, errors.New("разделитель не может содержать символы новой строки")
	}

	norm := normalizeRanges(parsed)

	return &Cut{
		cfg:    cfg,
		ranges: norm,
		delim:  []byte(cfg.Delimiter),
	}, nil
}

// parseFieldsSpec парсит строку спецификации полей
func parseFieldsSpec(spec string) ([]rangeSpec, error) {
	parts := strings.Split(spec, ",")

	var ranges []rangeSpec

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		if strings.Contains(p, "-") {
			sub := strings.SplitN(p, "-", 2)
			left := strings.TrimSpace(sub[0])
			right := strings.TrimSpace(sub[1])

			switch {
			case left == "" && right == "":
				return nil, fmt.Errorf("некорректный диапазон: %q", p)

			case left == "":
				r, err := strconv.Atoi(right)
				if err != nil || r < 1 {
					return nil, fmt.Errorf("некорректная верхняя граница диапазона %q", p)
				}
				ranges = append(ranges, rangeSpec{start: 1, end: r})

			case right == "":
				l, err := strconv.Atoi(left)
				if err != nil || l < 1 {
					return nil, fmt.Errorf("некорректная нижняя граница диапазона %q", p)
				}
				ranges = append(ranges, rangeSpec{start: l, end: -1})

			default:
				l, err1 := strconv.Atoi(left)
				r, err2 := strconv.Atoi(right)
				if err1 != nil || err2 != nil || l < 1 || r < 1 || l > r {
					return nil, fmt.Errorf("некорректный диапазон %q", p)
				}
				ranges = append(ranges, rangeSpec{start: l, end: r})
			}

		} else {
			n, err := strconv.Atoi(p)
			if err != nil || n < 1 {
				return nil, fmt.Errorf("некорректное число %q", p)
			}
			ranges = append(ranges, rangeSpec{start: n, end: n})
		}
	}

	return ranges, nil
}

// normalizeRanges сортирует диапазоны
func normalizeRanges(rs []rangeSpec) []rangeSpec {
	if len(rs) == 0 {
		return nil
	}

	rs = append([]rangeSpec(nil), rs...)

	sort.Slice(rs, func(i, j int) bool {
		if rs[i].start != rs[j].start {
			return rs[i].start < rs[j].start
		}
		ei := rs[i].end
		ej := rs[j].end
		if ei == -1 {
			return true
		}
		if ej == -1 {
			return false
		}
		return ei > ej
	})

	var out []rangeSpec
	cur := rs[0]

	for i := 1; i < len(rs); i++ {
		next := rs[i]

		if cur.end == -1 {
			break
		}

		if next.end == -1 {
			if next.start <= cur.end+1 {
				cur.end = -1
				break
			}
			out = append(out, cur)
			cur = next
			continue
		}

		if next.start <= cur.end+1 {
			if next.end > cur.end {
				cur.end = next.end
			}
		} else {
			out = append(out, cur)
			cur = next
		}
	}

	out = append(out, cur)
	return out
}
