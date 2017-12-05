package main

import (
	"fmt"
	"os"

	"github.com/justintout/advent-of-code/aocutil"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage:\n\t%s [input file]\n", os.Args[0])
		os.Exit(1)
	}
	ins := aocutil.ReadInput(os.Args[1])
	in := ins["input.txt"]
	fmt.Printf("%s\n", in)
}
