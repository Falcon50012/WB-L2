package sorter

import "testing"

func TestUnique(t *testing.T) {
	input := []Record{
		{Line: "a"},
		{Line: "a"},
		{Line: "b"},
	}

	result := Unique(input)

	if len(result) != 2 {
		t.Fatalf("expected 2 records, got %d", len(result))
	}

	if result[0].Line != "a" || result[1].Line != "b" {
		t.Fatal("unexpected unique result")
	}
}
