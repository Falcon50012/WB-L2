package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Shell struct {
	lastCode int
}

func New() *Shell { return &Shell{} }

func (s *Shell) Run(r io.Reader) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	defer signal.Stop(sigCh)

	go func() {
		for range sigCh {
			fmt.Fprintln(os.Stderr)
		}
	}()

	interactive := isTerminal(r)
	reader := bufio.NewReader(r)

	for {
		if interactive {
			cwd, _ := os.Getwd()
			fmt.Fprintf(os.Stderr, "%s$ ", cwd)
		}

		line, err := reader.ReadString('\n')
		line = strings.TrimRight(line, "\r\n")

		if err != nil {
			if err == io.EOF {
				if interactive {
					fmt.Fprintln(os.Stderr, "exit")
				}
				return
			}
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		cl := Parse(line)
		s.runList(cl)
	}
}

func isTerminal(r io.Reader) bool {
	f, ok := r.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
