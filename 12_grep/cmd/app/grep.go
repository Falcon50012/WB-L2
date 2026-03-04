package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Falcon50012/WB-L2/12_grep/internal/grep"
	"github.com/spf13/pflag"
)

func main() {
	cfg := grep.Config{}

	pflag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Использование: grep [ПАРАМЕТР]… ШАБЛОН [ФАЙЛ]…")
		pflag.PrintDefaults()
	}

	pflag.BoolVarP(&cfg.IgnoreCase, "ignore-case", "i", false, "игнорировать различие регистра")
	pflag.BoolVarP(&cfg.InvertMatch, "invert-match", "v", false, "выбирать не подходящие строки")
	pflag.BoolVarP(&cfg.LineNumber, "line-number", "n", false, "печатать номер строки")
	pflag.BoolVarP(&cfg.Fixed, "fixed-regexp", "F", false, "ШАБЛОН — строка, а не regexp")
	pflag.IntVarP(&cfg.After, "after-context", "A", 0, "печатать N строк после совпадения")
	pflag.IntVarP(&cfg.Before, "before-context", "B", 0, "печатать N строк перед совпадением")
	pflag.IntVarP(&cfg.Context, "context", "C", 0, "печатать N строк контекста")
	pflag.BoolVarP(&cfg.Count, "count", "c", false, "печатать только количество совпадений")

	pflag.Parse()

	if cfg.After < 0 || cfg.Before < 0 || cfg.Context < 0 {
		fmt.Fprintln(os.Stderr, "контекст (-A, -B, -C) не может быть отрицательным")
		os.Exit(2)
	}

	args := pflag.Args()
	if len(args) < 1 {
		pflag.Usage()
		os.Exit(2)
	}

	cfg.Pattern = args[0]

	var files []string
	if len(args) > 1 {
		files = args[1:]
	} else {
		files = []string{"-"}
	}

	g, err := grep.NewGrep(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	multiple := len(files) > 1
	totalMatches := 0
	hadOpenError := false

	for _, name := range files {
		var r io.Reader
		var f *os.File
		if name == "-" {
			r = os.Stdin
		} else {
			f, err = os.Open(name)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				hadOpenError = true
				continue
			}
			r = f
		}

		count, err := g.ProcessFile(os.Stdout, r, name, multiple)

		if f != nil {
			if cerr := f.Close(); cerr != nil {
				fmt.Fprintln(os.Stderr, "ошибка закрытия файла:", cerr)
			}
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}

		totalMatches += count
	}

	if hadOpenError {
		os.Exit(2)
	}

	if totalMatches == 0 {
		os.Exit(1)
	}
}
