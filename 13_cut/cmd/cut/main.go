package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Falcon50012/WB-L2/13_cut/internal/cut"
	"github.com/spf13/pflag"
)

func main() {
	cfg := cut.Config{}

	pflag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Использование: cut ПАРАМЕТР… [ФАЙЛ]…")
		pflag.PrintDefaults()
	}

	pflag.StringVarP(&cfg.Fields, "fields", "f", "", "выбрать только заданные поля")
	pflag.StringVarP(&cfg.Delimiter, "delimiter", "d", "\t", "использовать другой разделитель")
	pflag.BoolVarP(&cfg.OnlyDelimited, "only-delimited", "s", false, "не печатать строки без разделителя")

	pflag.Parse()

	args := pflag.Args()

	files := args
	if len(files) == 0 {
		files = []string{"-"}
	}

	c, err := cut.NewCut(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка конфигурации:", err)
		os.Exit(2)
	}

	var hadError bool

	for _, name := range files {
		var r io.Reader
		var f *os.File

		if name == "-" {
			r = os.Stdin
		} else {
			f, err = os.Open(name)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				hadError = true
				continue
			}
			r = f
		}

		if perr := c.ProcessFile(os.Stdout, r); perr != nil {
			fmt.Fprintln(os.Stderr, perr)
			hadError = true
		}

		if f != nil {
			if cerr := f.Close(); cerr != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при закрытии файла:", cerr)
				hadError = true
			}
		}
	}

	if hadError {
		os.Exit(2)
	}
}
