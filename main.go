package main

import (
	"fmt"
	"os"
)

func main() {
	if err := unescape(os.Stdin, os.Stdout); err != nil {
		os.Stdout.Write([]byte(fmt.Sprintf("%s: %s", os.Args[0], err.Error())))
	}
}
