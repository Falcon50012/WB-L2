package grep

import (
	"bytes"
	"strings"
	"testing"
)

func runGrep(t *testing.T, cfg Config, input, filename string, multiple bool) (string, int) {
	t.Helper()
	g, err := NewGrep(cfg)
	if err != nil {
		t.Fatalf("NewGrep error: %v", err)
	}
	var out bytes.Buffer
	cnt, err := g.ProcessFile(&out, strings.NewReader(input), filename, multiple)
	if err != nil {
		t.Fatalf("ProcessFile error: %v", err)
	}
	return out.String(), cnt
}

func TestSimpleMatch(t *testing.T) {
	cfg := Config{Pattern: "привет"}
	in := "hello\nпривет\nbye\n"
	out, cnt := runGrep(t, cfg, in, "-", false)
	if cnt != 1 {
		t.Fatalf("1 match expected, got %d", cnt)
	}
	if out != "привет\n" {
		t.Fatalf("unexpected out:\n%s", out)
	}
}

func TestIgnoreCase(t *testing.T) {
	cfg := Config{Pattern: "ПриВеТ", IgnoreCase: true}
	in := "привет\nПрИвЕт\nother\n"
	out, cnt := runGrep(t, cfg, in, "-", false)
	if cnt != 2 {
		t.Fatalf("2 match expected, got %d", cnt)
	}
	expected := "привет\nПрИвЕт\n"
	if out != expected {
		t.Fatalf("out not matched. Expect:\n%q\nGot:\n%q", expected, out)
	}
}

func TestInvertMatch(t *testing.T) {
	cfg := Config{Pattern: "да", InvertMatch: true}
	in := "да\nнет\nможет\n"
	out, cnt := runGrep(t, cfg, in, "-", false)
	if cnt != 2 {
		t.Fatalf("2 match expected, got %d", cnt)
	}
	expected := "нет\nможет\n"
	if out != expected {
		t.Fatalf("unexpected out:\n%q", out)
	}
}

func TestCount(t *testing.T) {
	cfg := Config{Pattern: "x", Count: true}
	in := "x\nx\n"
	out, cnt := runGrep(t, cfg, in, "-", false)
	if cnt != 2 {
		t.Fatalf("2 match expected, got %d", cnt)
	}
	if out != "2\n" {
		t.Fatalf("unexpected count: %q", out)
	}
}

func TestContext(t *testing.T) {
	cfg := Config{Pattern: "MATCH", Before: 1, LineNumber: true}
	in := "один\nдва MATCH\nтри\nчетыре\nпять MATCH\nшесть\n"
	out, cnt := runGrep(t, cfg, in, "-", false)
	if cnt != 2 {
		t.Fatalf("2 match expected, got %d", cnt)
	}
	expected := "1-один\n2:два MATCH\n--\n4-четыре\n5:пять MATCH\n"
	if out != expected {
		t.Fatalf("unexpected out. Expect:\n%q\nGot:\n%q", expected, out)
	}
}

func TestLastLine(t *testing.T) {
	cfg := Config{Pattern: "последняя"}
	in := "первая\nпоследняя строка"
	out, cnt := runGrep(t, cfg, in, "-", false)
	if cnt != 1 {
		t.Fatalf("1 match expected, got %d", cnt)
	}
	if out != "последняя строка\n" {
		t.Fatalf("unexpected out. Got: %q", out)
	}
}

func TestFixedVsRegex(t *testing.T) {
	cfg := Config{Pattern: "a.b", Fixed: true}
	in := "a.b\nacb\n"
	out, cnt := runGrep(t, cfg, in, "-", false)
	if cnt != 1 {
		t.Fatalf("1 match expected, got %d", cnt)
	}
	if out != "a.b\n" {
		t.Fatalf("unexpected out. Got: %q", out)
	}

	cfg = Config{Pattern: "a.b", Fixed: false}
	in = "a.b\nacb\n"
	out, cnt = runGrep(t, cfg, in, "-", false)
	if cnt != 2 {
		t.Fatalf("2 match expected, got %d", cnt)
	}
	if out != "a.b\nacb\n" {
		t.Fatalf("unexpected out. Got: %q", out)
	}
}
