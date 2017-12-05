package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage:\n\t%s [input file]\n", os.Args[0])
		os.Exit(1)
	}
	f, err := os.Open(os.Args[1])
	defer f.Close()
	if err != nil {
		log.Fatalf("unable to read %s: %v\n", os.Args[1], err)
	}
	reader := csv.NewReader(f)
	reader.Comma = '\t'
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("csv reader failed to read file: %v\n", err)
	}
	sum := 0
	diffs := make(chan int)
	var wg sync.WaitGroup
	for _, row := range data {
		wg.Add(1)
		go func(wg *sync.WaitGroup, diffs chan int, strings []string) {
			max := forceInt(strings[0])
			min := forceInt(strings[0])
			for _, s := range strings {
				i := forceInt(s)
				if i > max {
					max = i
				}
				if i < min {
					min = i
				}
			}
			diffs <- max - min
			wg.Done()
		}(&wg, diffs, row)
	}
	go func(wg *sync.WaitGroup, c chan int) {
		wg.Wait()
		close(c)
	}(&wg, diffs)
	for i := range diffs {
		sum += i
	}
	fmt.Printf("INPUT CHECKSUM: %d\n", sum)
}

func forceInt(s string) (i int) {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("unable to convert %s to int: %v", s, err)
	}
	return
}
