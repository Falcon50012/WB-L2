package sorter

import (
	"testing"
)

func TestCompareLexical(t *testing.T) {
	var opts Options

	a := Record{Key: "a"}
	b := Record{Key: "b"}

	compare := Compare(opts)

	if compare(a, b) >= 0 {
		t.Fatal("expected a < b")
	}
}

func TestCompareNumeric(t *testing.T) {
	opts := Options{Numeric: true}

	a := Record{Key: "10", Num: 10, NumOK: true}
	b := Record{Key: "2", Num: 2, NumOK: true}

	compare := Compare(opts)

	if compare(a, b) <= 0 {
		t.Fatal("expected 10 > 2")
	}
}

func TestFallbackToLexicalCompare(t *testing.T) {
	opts := Options{Numeric: true}

	a := Record{Key: "b", NumOK: false}
	b := Record{Key: "2", Num: 2, NumOK: true}

	compare := Compare(opts)

	if compare(a, b) <= 0 {
		t.Fatal("expected fallback to lexical compare")
	}
}

func Test_cmp(t *testing.T) {
	if cmp("a", "b") != -1 {
		t.Fatal("expected -1")
	}

	if cmp[float64](2, 1) != 1 {
		t.Fatal("expected 1")
	}
}
