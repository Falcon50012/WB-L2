package cut

import (
	"bytes"
	"strings"
	"testing"
)

func equalRanges(a, b []rangeSpec) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].start != b[i].start || a[i].end != b[i].end {
			return false
		}
	}
	return true
}

func TestParseFieldsSpec(t *testing.T) {
	cases := []struct {
		in  string
		exp []rangeSpec
	}{
		{"1", []rangeSpec{{1, 1}}},
		{"1,3-5", []rangeSpec{{1, 1}, {3, 5}}},
		{"-4", []rangeSpec{{1, 4}}},
		{"5-", []rangeSpec{{5, -1}}},
		{"1, 2-3,6-, 10", []rangeSpec{{1, 1}, {2, 3}, {6, -1}, {10, 10}}},
	}

	for _, c := range cases {
		got, err := parseFieldsSpec(c.in)
		if err != nil {
			t.Fatalf("parseFieldsSpec(%q) unexpected error: %v", c.in, err)
		}
		if !equalRanges(got, c.exp) {
			t.Fatalf("parseFieldsSpec(%q) = %v; want %v", c.in, got, c.exp)
		}
	}
}

func TestNormalizeRanges(t *testing.T) {
	in := []rangeSpec{
		{3, 5},
		{1, 1},
		{3, 10},
		{2, 6},
		{3, -1},
	}
	out := normalizeRanges(in)

	exp := []rangeSpec{{1, -1}}

	if !equalRanges(out, exp) {
		t.Fatalf("normalizeRanges(%v) = %v; want %v", in, out, exp)
	}
}

func TestJoinWithDelim(t *testing.T) {
	cases := []struct {
		fields [][]byte
		delim  []byte
		exp    []byte
	}{
		{[][]byte{}, []byte{'\t'}, []byte{}},
		{[][]byte{[]byte("a")}, []byte{'\t'}, []byte("a")},
		{[][]byte{[]byte("a"), []byte("b")}, []byte{','}, []byte("a,b")},
		{[][]byte{[]byte(""), []byte("")}, []byte{','}, []byte(",")},
	}

	for _, c := range cases {
		got := joinWithDelim(c.fields, c.delim)
		if !bytes.Equal(got, c.exp) {
			t.Fatalf("joinWithDelim(%v,%q) = %q; want %q", c.fields, c.delim, got, c.exp)
		}
	}
}

func runProcessFile(t *testing.T, cfg Config, input string) string {
	t.Helper()
	c, err := NewCut(cfg)
	if err != nil {
		t.Fatalf("NewCut(%v) error: %v", cfg, err)
	}
	var out bytes.Buffer
	if err := c.ProcessFile(&out, strings.NewReader(input)); err != nil {
		t.Fatalf("ProcessFile error: %v", err)
	}
	return out.String()
}

func TestProcessFile_Basic(t *testing.T) {
	cfg := Config{Fields: "1"}
	in := "a\tb\tc\nx\ty\tz\n"
	got := runProcessFile(t, cfg, in)
	want := "a\nx\n"
	if got != want {
		t.Fatalf("basic: got %q want %q", got, want)
	}
}

func TestProcessFile_MultipleFieldsAndRanges(t *testing.T) {
	cfg := Config{Fields: "1,3-4"}
	in := "a\tb\tc\td\n"
	got := runProcessFile(t, cfg, in)
	if got != "a\tc\td\n" {
		t.Fatalf("multi fields: got %q", got)
	}
}

func TestProcessFile_OnlyDelimited(t *testing.T) {
	cfg := Config{Fields: "1", OnlyDelimited: true}
	in := "no_delim_line\nwith:delim\n"
	got := runProcessFile(t, cfg, in)
	if got != "" {
		t.Fatalf("-s with no delim: got %q want empty", got)
	}

	cfg2 := Config{Fields: "1", Delimiter: ":", OnlyDelimited: true}
	in2 := "no_delim_line\nfirst:second\n"
	got2 := runProcessFile(t, cfg2, in2)
	if got2 != "first\n" {
		t.Fatalf("-s with colon delim: got %q want %q", got2, "first\n")
	}
}

func TestProcessFile_TrailingDelimiterEmptyField(t *testing.T) {
	cfg := Config{Fields: "3", Delimiter: ":"}
	in := "a:b:\n"
	got := runProcessFile(t, cfg, in)
	if got != "\n" {
		t.Fatalf("trailing delim: got %q want %q", got, "\n")
	}
}

func TestProcessFile_MultiByteDelimiter(t *testing.T) {
	cfg := Config{Fields: "2", Delimiter: "::"}
	in := "x::y::z\n"
	got := runProcessFile(t, cfg, in)
	if got != "y\n" {
		t.Fatalf("multibyte delim: got %q want %q", got, "y\n")
	}
}

func TestProcessFile_OpenEndedRange(t *testing.T) {
	cfg := Config{Fields: "2-", Delimiter: ":"}
	in := "a:b:c:d\n"
	got := runProcessFile(t, cfg, in)
	if got != "b:c:d\n" {
		t.Fatalf("open-ended range: got %q want %q", got, "b:c:d\n")
	}
}
