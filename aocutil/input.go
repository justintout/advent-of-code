package aocutil

import (
	"fmt"
	"io/ioutil"
	"os"
)

type fileResult struct {
	fn  string
	c   []byte
	err error
}

// ReadInput reads each of the files in fns and returns a map of the filenames to a []byte of their contents.
// Reads each input in parallel. Exits on error, since all inputs are assumed necessary
func ReadInput(fns ...string) map[string][]byte {
	o := make(map[string][]byte)

	res := make(chan *fileResult, len(fns))

	for _, fn := range fns {
		go func(res chan *fileResult, fn string) {
			b, err := ioutil.ReadFile(fn)
			r := &fileResult{
				fn:  fn,
				c:   b,
				err: err,
			}
			res <- r
		}(res, fn)
	}
	errState := false
	for i := 0; i < len(fns); i++ {
		r := <-res
		if r.err != nil {
			fmt.Fprintf(os.Stderr, "error reading input file %s: %v", r.fn, r.err)
			errState = true
			continue
		}
		if !errState {
			o[r.fn] = r.c
		}
	}
	if errState {
		os.Exit(1)
	}
	return o
}
