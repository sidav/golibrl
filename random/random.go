package random

import (
	"time"
)

// DEPRECATED: move to random_struct instead!

const (
	a = 2416
	c = 374441
	// m = 1771875
	m = 1000000005721 // prime
)

var (
	x int
)

func Randomize() int {
	x = int(time.Duration(time.Now().UnixNano())/time.Millisecond) % m
	return x
}

func SetSeed(val int) {
	x = val
}

func Random(modulo int) int {
	x = (x*a + c) % m
	if modulo != 0 {
		return x % modulo
	} else {
		return x
	}
}

func RollDice(dnum, dval, dmod int) int {
	var result int
	for i := 0; i < dnum; i++ {
		result += Random(dval) + 1
	}
	return result + dmod
}

func RandomUnitVectorInt() (int, int) {
	var vx, vy int
	for vx == 0 && vy == 0 {
		vx, vy = Random(3)-1, Random(3)-1
	}
	return vx, vy
}

func RandInRange(from, to int) int { //should be inclusive
	if to < from {
		t := from
		from = to
		to = t
	}
	if from == to {
		return from
	}
	return Random(to-from+1) + from
}

func RandomPercent() int {
	return Random(100)
}

func RandomCoordsInRangeFrom(x, y, r int) (int, int) {
	rx, ry := x+3*r, y+3*r
	for (rx-x)*(rx-x)+(ry-y)*(ry-y) > r*r {
		rx = RandInRange(x-r-1, x+r+1)
		ry = RandInRange(y-r-1, y+r+1)
	}
	return rx, ry
}

