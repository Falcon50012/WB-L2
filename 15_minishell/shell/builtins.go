package shell

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func isBuiltin(name string) bool {
	switch name {
	case "cd", "pwd", "echo", "kill", "ps":
		return true
	default:
		return false
	}
}

func runBuiltin(cmd *Command, stdout io.Writer, stderr io.Writer) int {
	switch cmd.Args[0] {
	case "cd":
		return builtinCD(cmd.Args[1:], stderr)
	case "pwd":
		return builtinPWD(stdout, stderr)
	case "echo":
		return builtinEcho(cmd.Args[1:], stdout)
	case "kill":
		return builtinKill(cmd.Args[1:], stderr)
	case "ps":
		return builtinPS(stdout, stderr)
	default:
		return 1
	}
}

func builtinCD(args []string, stderr io.Writer) int {
	dir := ""
	if len(args) == 0 {
		dir = os.Getenv("HOME")
		if dir == "" {
			fmt.Fprintln(stderr, "cd: HOME not set")
			return 1
		}
	} else {
		dir = args[0]
	}

	if dir == "~" {
		dir = os.Getenv("HOME")
	} else if strings.HasPrefix(dir, "~/") {
		dir = filepath.Join(os.Getenv("HOME"), dir[2:])
	}

	if err := os.Chdir(dir); err != nil {
		fmt.Fprintf(stderr, "cd: %v\n", err)
		return 1
	}
	return 0
}

func builtinPWD(stdout io.Writer, stderr io.Writer) int {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(stderr, "pwd: %v\n", err)
		return 1
	}
	fmt.Fprintln(stdout, dir)
	return 0
}

func builtinEcho(args []string, stdout io.Writer) int {
	noNewline := false
	start := 0
	if len(args) > 0 && args[0] == "-n" {
		noNewline = true
		start = 1
	}

	out := strings.Join(args[start:], " ")
	if noNewline {
		fmt.Fprint(stdout, out)
	} else {
		fmt.Fprintln(stdout, out)
	}
	return 0
}

func builtinKill(args []string, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "kill: usage: kill [-signal] pid...")
		return 1
	}

	sig := syscall.SIGTERM
	pidArgs := args

	if strings.HasPrefix(args[0], "-") {
		spec := args[0][1:]
		var err error
		sig, err = parseSignal(spec)
		if err != nil {
			fmt.Fprintf(stderr, "kill: %v\n", err)
			return 1
		}
		pidArgs = args[1:]
	}

	if len(pidArgs) == 0 {
		fmt.Fprintln(stderr, "kill: missing pid")
		return 1
	}

	code := 0
	for _, s := range pidArgs {
		pid, err := strconv.Atoi(s)
		if err != nil || pid <= 0 {
			fmt.Fprintf(stderr, "kill: invalid pid: %q\n", s)
			code = 1
			continue
		}
		if err := syscall.Kill(pid, sig); err != nil {
			fmt.Fprintf(stderr, "kill: %v\n", err)
			code = 1
		}
	}
	return code
}

func parseSignal(s string) (syscall.Signal, error) {
	upper := strings.ToUpper(s)
	upper = strings.TrimPrefix(upper, "SIG")

	switch upper {
	case "2", "INT":
		return syscall.SIGINT, nil
	case "9", "KILL":
		return syscall.SIGKILL, nil
	case "15", "TERM":
		return syscall.SIGTERM, nil
	default:
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("unknown signal: %q", s)
		}
		return syscall.Signal(n), nil
	}
}

func builtinPS(stdout io.Writer, stderr io.Writer) int {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		fmt.Fprintf(stderr, "ps: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "%6s  %s\n", "PID", "COMMAND")
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		commBytes, err := os.ReadFile("/proc/" + entry.Name() + "/comm")
		if err != nil {
			continue
		}
		comm := strings.TrimSpace(string(commBytes))

		fmt.Fprintf(stdout, "%6d  %s\n", pid, comm)
	}
	return 0
}
