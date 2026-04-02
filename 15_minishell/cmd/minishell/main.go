package main

import (
	"os"

	"github.com/Falcon50012/WB-L2/15/shell"
)

func main() {
	sh := shell.New()
	sh.Run(os.Stdin)
}
