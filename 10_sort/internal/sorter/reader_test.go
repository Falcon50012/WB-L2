package sorter

import (
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	line := "a\t1\nb\t2\n"

	tmp, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.WriteString(line); err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	opts := Options{
		Column:  2,
		Numeric: true,
	}

	records, err := Read(tmp.Name(), opts)
	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}

	if !records[0].NumOK {
		t.Fatal("expected numeric parsing")
	}
}
