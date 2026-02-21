package sorter

import "testing"

func TestSortLexical(t *testing.T) {
	var opts Options

	input := []Record{
		{Key: "b"},
		{Key: "a"},
	}

	Sort(input, opts)

	if input[0].Key != "a" {
		t.Fatal("expected sorted order")
	}
}

func TestSortNumeric(t *testing.T) {
	opts := Options{Numeric: true}

	input := []Record{
		{Key: "10", Num: 10, NumOK: true},
		{Key: "2", Num: 2, NumOK: true},
	}

	Sort(input, opts)

	if input[0].Num != 2 {
		t.Fatal("expected numeric sorting")
	}
}
