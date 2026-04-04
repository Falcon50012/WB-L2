package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Falcon50012/WB-L2/16_wgetm/internal/wgetm"
)

func main() {
	depth := flag.Int("d", 2, "recuesion depth")
	workers := flag.Int("w", 5, "workers num")
	output := flag.String("o", ".", "output dir")
	timeout := flag.Duration("t", 30*time.Second, "request timeout")
	verbose := flag.Bool("v", false, "verbose")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Using: wget [flags] <URL>")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	cfg := wgetm.Config{
		StartURL:  flag.Arg(0),
		MaxDepth:  *depth,
		Workers:   *workers,
		OutputDir: *output,
		Timeout:   *timeout,
		Verbose:   *verbose,
	}

	c, err := wgetm.NewCrawler(cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	stats := c.Run()
	fmt.Printf("Download: %d  Skipped: %d  Errors: %d\n",
		stats.Downloaded, stats.Skipped, stats.Errors)
}
