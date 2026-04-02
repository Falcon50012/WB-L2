package shell

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func (s *Shell) runList(cl *CommandList) {
	if cl == nil || len(cl.Items) == 0 {
		return
	}

	s.lastCode = s.runPipeline(cl.Items[0].Pipeline)

	for _, item := range cl.Items[1:] {
		switch item.Op {
		case "&&":
			if s.lastCode == 0 {
				s.lastCode = s.runPipeline(item.Pipeline)
			}
		case "||":
			if s.lastCode != 0 {
				s.lastCode = s.runPipeline(item.Pipeline)
			}
		}
	}
}

func (s *Shell) runPipeline(pl *Pipeline) int {
	n := len(pl.Cmds)
	if n == 0 {
		return 0
	}

	if n == 1 {
		return s.runCmd(pl.Cmds[0], os.Stdin, os.Stdout)
	}

	rds := make([]*os.File, n-1)
	wrs := make([]*os.File, n-1)

	for i := range rds {
		r, w, err := os.Pipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "minishell: pipe: %v\n", err)
			for j := range i {
				_ = rds[j].Close()
				_ = wrs[j].Close()
			}
			return 1
		}
		rds[i], wrs[i] = r, w
	}

	codes := make([]int, n)
	var wg sync.WaitGroup

	for i, cmd := range pl.Cmds {
		var in, out *os.File
		if i == 0 {
			in = os.Stdin
		} else {
			in = rds[i-1]
		}
		if i == n-1 {
			out = os.Stdout
		} else {
			out = wrs[i]
		}

		wg.Add(1)
		go func(idx int, c *Command, in, out *os.File) {
			defer wg.Done()
			defer closePipe(in)
			defer closePipe(out)

			if c.Args[0] == "cd" {
				fmt.Fprintln(os.Stderr, "minishell: cd: cannot be used in a pipeline")
				codes[idx] = 1
				return
			}

			codes[idx] = s.runCmd(c, in, out)
		}(i, cmd, in, out)
	}

	wg.Wait()
	return codes[n-1]
}

func closePipe(f *os.File) {
	if f != nil && f != os.Stdin && f != os.Stdout {
		_ = f.Close()
	}
}

func (s *Shell) runCmd(cmd *Command, stdin, stdout *os.File) int {
	if len(cmd.Args) == 0 {
		return 0
	}

	if cmd.RedirIn != "" {
		f, err := os.Open(cmd.RedirIn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "minishell: %v\n", err)
			return 1
		}
		defer f.Close()
		stdin = f
	}

	if cmd.RedirOut != "" {
		f, err := os.OpenFile(cmd.RedirOut, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "minishell: %v\n", err)
			return 1
		}
		defer f.Close()
		stdout = f
	}

	if isBuiltin(cmd.Args[0]) {
		return runBuiltin(cmd, stdout, os.Stderr)
	}

	return s.runExternal(cmd, stdin, stdout)
}

func (s *Shell) runExternal(cmd *Command, stdin, stdout *os.File) int {
	c := exec.Command(cmd.Args[0], cmd.Args[1:]...)
	c.Stdin = stdin
	c.Stdout = stdout
	c.Stderr = os.Stderr

	if err := c.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "minishell: %s: %v\n", cmd.Args[0], err)
		return 127
	}

	if err := c.Wait(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return ee.ExitCode()
		}
		return 1
	}

	return 0
}
