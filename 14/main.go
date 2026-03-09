package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

func or(channels ...<-chan any) <-chan any {
	done := make(chan any)
	out := orRecursive(done, channels...)

	go func() {
		<-out
		close(done)
	}()

	return out
}

func orRecursive(done <-chan any, channels ...<-chan any) <-chan any {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	out := make(chan any)

	go func() {
		defer close(out)

		mid := len(channels) / 2

		select {
		case <-done:
		case <-orRecursive(done, channels[:mid]...):
		case <-orRecursive(done, channels[mid:]...):
		}
	}()

	return out
}

func main() {
	sig := func(after time.Duration) <-chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v\n", time.Since(start))

	p := pprof.Lookup("goroutineleak")
	if p != nil {
		p.WriteTo(os.Stdout, 1)
	} else {
		fmt.Println("goroutine leak profile not available (requires GOEXPERIMENT=goroutineleakprofile)")
	}
}
