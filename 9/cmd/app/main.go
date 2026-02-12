package main

import (
	"fmt"
	"log"

	"github.com/Falcon50012/WB-L2/9/unpack"
)

func main() {
	s := "\\"
	unpacked, err := unpack.Unpack(s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("unpacked string:", unpacked)
}
