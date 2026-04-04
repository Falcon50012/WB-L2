package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: telnet [--timeout=<duration>s] <host> <port>")
		os.Exit(1)
	}

	host, port := args[0], args[1]
	address := net.JoinHostPort(host, port)

	conn, err := net.DialTimeout("tcp", address, *timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection failed: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Fprintf(os.Stderr, "Connected to %s\n", address)

	done := make(chan struct{})
	var once sync.Once
	finish := func() { once.Do(func() { close(done) }) }

	go func() {
		defer finish()

		r := bufio.NewReader(os.Stdin)
		w := bufio.NewWriter(conn)

		for {
			line, err := r.ReadString('\n')

			if len(line) > 0 {
				if _, werr := w.WriteString(line); werr != nil {
					fmt.Fprintf(os.Stderr, "Write error: %v\n", werr)
					return
				}
				if werr := w.Flush(); werr != nil {
					fmt.Fprintf(os.Stderr, "Flush error: %v\n", werr)
					return
				}
			}

			if err != nil {
				if err == io.EOF {
					fmt.Fprintln(os.Stderr, "EOF")
					if tc, ok := conn.(*net.TCPConn); ok {
						tc.CloseWrite()
					}
				} else {
					fmt.Fprintf(os.Stderr, "Stdin error: %v\n", err)
				}
				return
			}
		}
	}()

	go func() {
		defer finish()

		s := bufio.NewScanner(conn)
		for s.Scan() {
			fmt.Println(s.Text())
		}
		if err := s.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
		} else {
			fmt.Fprintln(os.Stderr, "Connection closed by server")
		}
	}()

	<-done
}
