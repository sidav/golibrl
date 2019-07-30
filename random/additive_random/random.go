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

func (f *FibRandom) lcg() int { // used for initial values setup
	f.lcgX = (f.lcgX*2416 + 374441) % 1771875
	return f.lcgX
}

func (f *FibRandom) InitDefault() {
	f.InitCustom(-1, 17, 5)
}

func (f *FibRandom) InitBySeed(seed int) {
	f.InitCustom(seed, 17, 5)
}

func (f *FibRandom) InitCustom(seed, lagA, lagB int) {
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
	f.lcgX = seed
	f.biggerLag = lagA
	f.smallerLag = lagB
	f.curValues = make([]int, 0)
	for i := 0; i < f.biggerLag; i++ {
		newval := f.lcg() % mod
		f.curValues = append(f.curValues, newval)
	}
}

func (f *FibRandom) Rand(modulo int) int {
	aIndex := f.currIndex
	bIndex := f.currIndex - f.smallerLag
	if bIndex < 0 {
		bIndex += f.biggerLag
	}
	b := f.curValues[bIndex]
	a := f.curValues[aIndex]
	new := a + b
	if new >= mod {
		new -= mod
	}
	f.curValues[f.currIndex] = new
	f.currIndex++
	if f.currIndex >= len(f.curValues) {
		f.currIndex = 0
	}
	if modulo > 0 {
		return new % modulo
	}
	return new
}
