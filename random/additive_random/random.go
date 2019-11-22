package additive_random

import (
	"time"
)

const mod = 1<<31 - 1

type FibRandom struct {
	curValues             []int
	currIndex, lcgX       int
	biggerLag, smallerLag int // should be > 0
}

func (rnd *FibRandom) lcg() int { // used for initial values setup
	rnd.lcgX = (rnd.lcgX*2416 + 374441) % 1771875
	return rnd.lcgX
}

func (rnd *FibRandom) InitDefault() {
	rnd.InitCustom(-1, 17, 5)
}

func (rnd *FibRandom) InitBySeed(seed int) {
	rnd.InitCustom(seed, 17, 5)
}

func (rnd *FibRandom) InitCustom(seed, lagA, lagB int) {
	if seed < 0 {
		seed = int(time.Duration(time.Now().UnixNano())/time.Millisecond) % mod
	}
	if lagB > lagA {
		t := lagB
		lagB = lagA
		lagA = t
	}
	if lagB <= 0 || lagB == lagA {
		panic("FibRand lag params should be > 0 and not equal!")
	}
	rnd.lcgX = seed
	rnd.biggerLag = lagA
	rnd.smallerLag = lagB
	rnd.curValues = make([]int, 0)
	for i := 0; i < rnd.biggerLag; i++ {
		newval := rnd.lcg() % mod
		rnd.curValues = append(rnd.curValues, newval)
	}
}

func (rnd *FibRandom) Rand(modulo int) int {
	aIndex := rnd.currIndex
	bIndex := rnd.currIndex - rnd.smallerLag
	if bIndex < 0 {
		bIndex += rnd.biggerLag
	}
	b := rnd.curValues[bIndex]
	a := rnd.curValues[aIndex]
	new := a + b
	if new >= mod {
		new -= mod
	}
	rnd.curValues[rnd.currIndex] = new
	rnd.currIndex++
	if rnd.currIndex >= len(rnd.curValues) {
		rnd.currIndex = 0
	}
	if modulo > 0 {
		return new % modulo
	}
	return new
}

func (rnd *FibRandom) OneChanceFrom(numChances int) bool {
	return rnd.Rand(numChances) == 0 
}

func (rnd *FibRandom) BiasedRandInRange(from, to, bias, influencePercent int) int {
	// rnd = random() x (max - min) + min
	// mix = random() x influence
	// value = rnd x (1 - mix) + bias x mix
	const factor = 1
	influencePercent *= factor 
	totalrange := to - from
	rand := rnd.RandInRange(0, totalrange)
	mix := rnd.RandInRange(0, influencePercent)
	result := (rand * (100*factor - mix) + (bias-from) * mix)
	// proper rounding:  
	if result % (100*factor) >= (50*factor) {
		result += 100*factor 
	}
	result /= 100*factor

	return result + from 
}

func (rnd *FibRandom) RollDice(dnum, dval, dmod int) int {
	var result int
	for i := 0; i < dnum; i++ {
		result += rnd.Rand(dval) + 1
	}
	return result + dmod
}

func (rnd *FibRandom) RandomUnitVectorInt() (int, int) {
	var vx, vy int
	for vx == 0 && vy == 0 {
		vx, vy = rnd.Rand(3)-1, rnd.Rand(3)-1
	}
	return vx, vy
}

func (rnd *FibRandom) RandInRange(from, to int) int { //should be inclusive
	if to < from {
		t := from
		from = to
		to = t
	}
	if from == to {
		return from
	}
	return rnd.Rand(to-from+1) + from
}

func (rnd *FibRandom) RandomPercent() int {
	return rnd.Rand(100)
}

func (rnd *FibRandom) RandomCoordsInRangeFrom(x, y, r int) (int, int) {
	rx, ry := x+3*r, y+3*r
	for (rx-x)*(rx-x)+(ry-y)*(ry-y) > r*r {
		rx = rnd.RandInRange(x-r-1, x+r+1)
		ry = rnd.RandInRange(y-r-1, y+r+1)
	}
	return rx, ry
}

