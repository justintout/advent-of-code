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

	fmt.Printf("FILE CHECKSUM: %d\n", partOne(data))
	fmt.Printf("DIVISION SUM: %d\n", partTwo(data))

}

func forceInt(s string) (i int) {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("unable to convert %s to int: %v", s, err)
	}
	return
}

func partOne(data [][]string) int {
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
	return sum
}

func partTwo(data [][]string) int {
	prods := make(chan int)
	var wg sync.WaitGroup
	for _, row := range data {
		// fmt.Println(row)
		go divideRow(prods, &wg, row)
		wg.Add(1)
	}

	go func(c chan int, wg *sync.WaitGroup) {
		wg.Wait()
		close(c)
	}(prods, &wg)

	sum := 0
	for i := range prods {
		sum += i
	}
	return sum
}

func divideRow(res chan int, wg *sync.WaitGroup, row []string) {
	defer wg.Done()
	c := make(chan int)
	for i := 0; i < len(row); i++ {
		var r []string
		r = append(r, row[:i]...)
		r = append(r, row[i+1:]...)
		go func(c chan int, e string, r []string) {
			i := forceInt(e)
			for _, re := range r {
				if i%forceInt(re) == 0 {
					c <- i / forceInt(re)
					close(c)
				}
			}
		}(c, row[i], r)
	}
	result := <-c
	res <- result
}
