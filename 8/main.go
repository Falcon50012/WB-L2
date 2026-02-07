package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

const ntpServer = "0.ru.pool.ntp.org"

func main() {
	t, err := ntp.Time(ntpServer)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(t.Format(time.RFC3339Nano))
}
