package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Falcon50012/WB-L2/10_sort/internal/sorter"
)

// main — точка входа в программу
func main() {
	var opts sorter.Options

	flag.IntVar(&opts.Column, "k", 0, "sort by column number (1-based). 0 => whole line")
	flag.StringVar(&opts.Delimiter, "d", "\t", "set variable delimeter")
	flag.BoolVar(&opts.Numeric, "n", false, "numeric sort")
	flag.BoolVar(&opts.Reverse, "r", false, "reverse order")
	flag.BoolVar(&opts.Unique, "u", false, "output unique lines")
	flag.Parse()

	var filename string
	if flag.NArg() > 0 {
		filename = flag.Arg(0)
	}

	records, err := sorter.Read(filename, opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read error:", err)
		os.Exit(1)
	}

	sorter.Sort(records, opts)

	if opts.Unique {
		records = sorter.Unique(records)
	}

	if err := sorter.Write(records); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
