package random

import "time"

type LCGRandom struct {
	x int
}

func (r *LCGRandom) Randomize() int {
	r.x = int(time.Duration(time.Now().UnixNano())/time.Millisecond) % m
	return r.x
}

func  (r *LCGRandom) SetSeed(val int) {
	r.x = val
}

func (r *LCGRandom) Random(modulo int) int {
	r.x = (r.x*a + c) % m
	if modulo != 0 {
		return r.x % modulo
	} else {
		return r.x
	}
}

func (r *LCGRandom) RollDice(dnum, dval, dmod int) int {
	var result int
	for i := 0; i < dnum; i++ {
		result += r.Random(dval) + 1
	}
	return result + dmod
}

func (r *LCGRandom) RandomUnitVectorInt() (int, int) {
	var vx, vy int
	for vx == 0 && vy == 0 {
		vx, vy = r.Random(3)-1, r.Random(3)-1
	}
	return vx, vy
}

func (r *LCGRandom) RandInRange(from, to int) int { //should be inclusive
	if to < from {
		t := from
		from = to
		to = t
	}
	if from == to {
		return from
	}
	return r.Random(to-from+1) + from
}

func (r *LCGRandom) RandomPercent() int {
	return r.Random(100)
}

func (rnd *LCGRandom) RandomCoordsInRangeFrom(x, y, r int) (int, int) {
	rx, ry := x+3*r, y+3*r
	for (rx-x)*(rx-x)+(ry-y)*(ry-y) > r*r {
		rx = rnd.RandInRange(x-r-1, x+r+1)
		ry = rnd.RandInRange(y-r-1, y+r+1)
	}
	return rx, ry
}
