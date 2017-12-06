package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

const input = 325489

const (
	RIGHT = iota
	TOP
	LEFT
	BOTTOM
	ERROR
)

type ring struct {
	n, br, tr, tl, bl, l, min, max, rc, tc, lc, bc int
}

func newRing(n int) *ring {
	l := 2*n - 1
	min := int(math.Pow(float64(l-2), 2) + 1)
	max := int(math.Pow(float64(l), 2))
	tr := min + l - 2
	tl := min + ((l - 1) * 2) - 1
	bl := max - l + 1
	return &ring{
		n:   n,
		l:   l,
		br:  max,
		tr:  tr,
		tl:  tl,
		bl:  bl,
		rc:  tr - l/2,
		tc:  tr + l/2,
		lc:  tl + l/2,
		bc:  max - l/2,
		min: min,
		max: max,
	}
}

func (r ring) String() string {
	return fmt.Sprintf("%d %d %d\n%d - %d\n%d %d %d\n", r.tl, r.tc, r.tr, r.lc, r.rc, r.bl, r.bc, r.br)
}

func (r ring) side(n int) int {
	if !within(n, r.min, r.max) {
		return ERROR
	}
	if within(n, r.tr, r.tl) {
		return TOP
	}
	if within(n, r.tl, r.bl) {
		return LEFT
	}
	if within(n, r.bl, r.max) {
		return BOTTOM
	}
	return RIGHT
}

func (r ring) distToSideCenter(n int) int {
	s := r.side(n)
	if s == ERROR {
		log.Fatalf("%d is not contained in ring %d", n, r.n)
	}
	d := 0
	switch s {
	case RIGHT:
		d = abs(n - r.rc)
	case TOP:
		d = abs(n - r.tc)
	case LEFT:
		d = abs(n - r.lc)
	case BOTTOM:
		d = abs(n - r.bc)
	}
	return d
}

func (r ring) distToOrigin(n int) int {
	return r.n - 1 + r.distToSideCenter(n)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage:\n\t%s [input]\n", os.Args[0])
		os.Exit(1)
	}
	input, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("input must be an integer\n")
	}
	fmt.Printf("DISTANCE TO ORIGIN: %d", partOne(input))
}

func partOne(input int) int {
	if input == 1 {
		return 0
	}
	r := newRing(1)
	for i := 2; !within(input, r.min, r.max); i++ {
		r = newRing(i)
	}
	return r.distToOrigin(input)
}

func within(n, l, u int) bool {
	return u >= n && n >= l
}

func abs(n int) int {
	if n < 0 {
		n = -n
	}
	return n
}

func square(n int) int {
	return int(math.Pow(float64(n), 2))
}
