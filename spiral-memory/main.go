package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

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
	fmt.Printf("DISTANCE TO ORIGIN: %d\n", partOne(input))
	fmt.Printf("NEXT VALUE: %d\n", partTwo(input))
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

func partTwo(input int) int {
	s := newSpiral()
	return s.generateUntil(input)
}

type spiral map[point]int

type direction byte

const (
	EAST direction = iota
	NORTH
	WEST
	SOUTH
)

func (d direction) String() string {
	s := ""
	switch d {
	case EAST:
		s = "east"
	case NORTH:
		s = "north"
	case WEST:
		s = "west"
	case SOUTH:
		s = "south"
	}
	return s
}

func newSpiral() spiral {
	m := map[point]int{
		point{0, 0}: 1,
	}
	s := spiral(m)
	return s
}

// strategy: we start a spiral with the center filled. we know the next is 0,1 = 1, so this is the starting condition.
// from there, we peek ahead to the left. if it already has a value we continue ahead. if it doesnt, we turn left and continue.
// at each step, we sum the neighbor values and set the current equal to this sum. we return the most recent value added. this
// will be the first value greater than or equal to the input.
func (s spiral) generateUntil(input int) int {
	p := point{1, 0}
	d := EAST
	s[p] = 1
	for i := 1; s[p] <= input; i++ {
		_, ok := s[p.ahead(turnLeft(d))]
		if !ok {
			d = turnLeft(d)
		}
		p = p.ahead(d)
		v := s.sumNeigbors(p)
		s[p] = v
		fmt.Printf("(%d, %d) = %d, %s\n", p.x, p.y, v, d)
	}
	return s[p]
}

func (s spiral) sumNeigbors(p point) int {
	t := 0
	for _, n := range p.neighbors() {
		v, ok := s[n]
		if ok {
			fmt.Printf("  (%d, %d) = %d\n", n.x, n.y, v)
			t += v
		}
	}
	return t

}

func turnLeft(d direction) direction {
	switch d {
	case EAST:
		return NORTH
	case NORTH:
		return WEST
	case WEST:
		return SOUTH
	}
	return EAST
}

type point struct {
	x, y int
}

func (p point) neighbors() []point {
	return []point{point{p.x - 1, p.y}, point{p.x - 1, p.y - 1}, point{p.x - 1, p.y + 1}, point{p.x, p.y + 1}, point{p.x, p.y - 1}, point{p.x + 1, p.y}, point{p.x + 1, p.y - 1}, point{p.x + 1, p.y + 1}}
}

func (p point) ahead(d direction) point {
	var a point
	switch d {
	case EAST:
		a.x = p.x + 1
		a.y = p.y
	case NORTH:
		a.x = p.x
		a.y = p.y + 1
	case WEST:
		a.x = p.x - 1
		a.y = p.y
	case SOUTH:
		a.x = p.x
		a.y = p.y - 1
	}
	return a
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
