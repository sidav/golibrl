package Fractal_landscape

import (
	"github.com/sidav/golibrl/random"
)

var rnd random.LCGRandom

func GenHeightMap(w, h int) *[][]int {
	rnd.Randomize()
	if w%2 == 0 {
		w++
	}
	if h%2 == 0 {
		h++
	}
	hmap := make([][]int, w)
	for i := range hmap {
		hmap[i] = make([]int, h)
	}
	initHeightMap(&hmap)
	iterate(&hmap)
	return &hmap
}

func initHeightMap(fullMap *[][]int) {
	w, h := len(*fullMap), len((*fullMap)[0])
	var (
		cornerFrom = -40
		cornerTo = -5
		centerFrom = 0
		centerTo   = 100
		subcFrom   = -20
		subcTo     = 65
	)
	// corners
	(*fullMap)[0][0] = rnd.RandInRange(cornerFrom, cornerTo)
	(*fullMap)[0][h-1] = rnd.RandInRange(cornerFrom, cornerTo)
	(*fullMap)[w-1][0] = rnd.RandInRange(cornerFrom, cornerTo)
	(*fullMap)[w-1][h-1] = rnd.RandInRange(cornerFrom, cornerTo)
	// center
	(*fullMap)[w/2][h/2] = rnd.RandInRange(centerFrom, centerTo)
	// fractal sub-centers
	(*fullMap)[w/4][h/4] = rnd.RandInRange(subcFrom, subcTo)
	(*fullMap)[w/4][h/2] = rnd.RandInRange(subcFrom, subcTo)
	(*fullMap)[w/4][h/2+h/4] = rnd.RandInRange(subcFrom, subcTo)
	(*fullMap)[w/2][h/4] = rnd.RandInRange(subcFrom, subcTo)
	(*fullMap)[w/2][h/2] = rnd.RandInRange(subcFrom, subcTo)
	(*fullMap)[w/2][h/2+h/4] = rnd.RandInRange(subcFrom, subcTo)
	(*fullMap)[w/2+w/4][h/4] = rnd.RandInRange(subcFrom, subcTo)
	(*fullMap)[w/2+w/4][h/2] = rnd.RandInRange(subcFrom, subcTo)
	(*fullMap)[w/2+w/4][h/2+h/4] = rnd.RandInRange(subcFrom, subcTo)
}

func iterate(fullMap *[][]int) {
	w, h := len(*fullMap), len((*fullMap)[0])
	sqw := w
	if h < w {
		sqw = h
	}
	for sqw > 2 {
		for x:=0;x<=w-sqw;x+=sqw-1 {
			for y:=0; y<=h-sqw;y+=sqw-1 {
				doSquareMidpoint(fullMap, x, y, sqw, -2, -2)
			}
		}
		sqw = sqw/2 + 1
	}
}

func doSquareMidpoint(fullMap *[][]int, x, y, w, spreadBorder, spreadCenter int) {
	if x < 0 || y < 0 || x+w > len(*fullMap) || y+w > len(*fullMap) {
		return
	}
	midX := x + w/2
	rightX := x + w - 1
	midY := y + w/2
	botY := y + w - 1
	sizeIsEven := w%2 == 0 
	// midpoint left side of square
	if (*fullMap)[x][midY] == 0 {
		(*fullMap)[x][midY] = jitterAvg2((*fullMap)[x][y], (*fullMap)[x][botY], spreadBorder)
	}
	// midpoint right side of square
	if (*fullMap)[rightX][midY] == 0 {
		(*fullMap)[rightX][midY] = jitterAvg2((*fullMap)[rightX][y], (*fullMap)[rightX][botY], spreadBorder)
	}
	// midpoint top side of a square
	if (*fullMap)[midX][y] == 0 {
		(*fullMap)[midX][y] = jitterAvg2((*fullMap)[x][y], (*fullMap)[rightX][y], spreadBorder)
	}
	// midpoint bottom side of a square
	if (*fullMap)[midX][botY] == 0 {
		(*fullMap)[midX][botY] = jitterAvg2((*fullMap)[x][botY], (*fullMap)[rightX][botY], spreadBorder)
	}
	//midpoint center of a square
	if (*fullMap)[midX][midY] == 0 {
		(*fullMap)[midX][midY] = jitterAvg4((*fullMap)[x][midY], (*fullMap)[rightX][midY], (*fullMap)[midX][y], (*fullMap)[midX][botY], spreadCenter)
	}
	if sizeIsEven {
		// midpoint left side of square
		if (*fullMap)[x][midY-1] == 0 {
			(*fullMap)[x][midY-1] = jitterAvg2((*fullMap)[x][y], (*fullMap)[x][botY], spreadBorder)
		}
		// midpoint right side of square
		if (*fullMap)[rightX][midY-1] == 0 {
			(*fullMap)[rightX][midY-1] = jitterAvg2((*fullMap)[rightX][y], (*fullMap)[rightX][botY], spreadBorder)
		}
		// midpoint top side of a square
		if (*fullMap)[midX-1][y] == 0 {
			(*fullMap)[midX-1][y] = jitterAvg2((*fullMap)[x][y], (*fullMap)[rightX][y], spreadBorder)
		}
		// midpoint bottom side of a square
		if (*fullMap)[midX-1][botY] == 0 {
			(*fullMap)[midX-1][botY] = jitterAvg2((*fullMap)[x][botY], (*fullMap)[rightX][botY], spreadBorder)
		}
		//midpoint center of a square
		(*fullMap)[midX-1][midY-1] = (*fullMap)[midX][midY]
		(*fullMap)[midX][midY-1] = (*fullMap)[midX][midY]
		(*fullMap)[midX-1][midY] = (*fullMap)[midX][midY]
	}
}

// boredom below

func isPowerOfTwo(x int) bool {
	return (x & (x - 1)) == 0
}

func sliceMinMax(arr []int) (int, int) {
	min := arr[0]
	max := arr[0]
	for i := range arr {
		if arr[i] > max {
			max = arr[i]
		}
		if arr[i] < min {
			min = arr[i]
		}
	}
	return min, max
}

func jitterAvg2(a1, a2, spread int) int {
	min, max := sliceMinMax([]int{a1, a2})
	return rnd.RandInRange(min-spread, max+spread)
}

func jitterAvg4(a1, a2, a3, a4, spread int) int {
	min, max := sliceMinMax([]int{a1, a2, a3, a4})
	return rnd.RandInRange(min-spread, max+spread)
}
