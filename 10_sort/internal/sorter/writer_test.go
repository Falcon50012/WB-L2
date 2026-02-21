package sorter

import (
	"bytes"
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	records := []Record{
		{Line: "a"},
		{Line: "b"},
	}

	var buf bytes.Buffer

	out := os.Stdout
	reader, writer, _ := os.Pipe()
	os.Stdout = writer

	Write(records)

	writer.Close()
	os.Stdout = out

	buf.ReadFrom(reader)

	expected := "a\nb\n"
	if buf.String() != expected {
		t.Fatalf("got %q, want %q", buf.String(), expected)
	}
}
