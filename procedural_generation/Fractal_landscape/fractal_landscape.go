package Fractal_landscape

import (
	"github.com/sidav/golibrl/random"
)

var rnd random.LCGRandom

//func avg2(a1, a2 int) int {
//	return (a1 + a2) / 2
//}
//
//func avg4(a1, a2, a3, a4 int) int {
//	return (a1 + a2 + a3 + a4) / 4
//}
//
//func jitter(value, spread int) int {
//	return rnd.RandInRange(value-spread, value+spread)
//}

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
	(*fullMap)[0][0] = -10
	(*fullMap)[0][h-1] = -10
	(*fullMap)[w-1][0] = -10
	(*fullMap)[w-1][h-1] = -10
	(*fullMap)[w/2][h/2] = 55
}

func iterate(fullMap *[][]int) {
	w, h := len(*fullMap), len((*fullMap)[0])
	sqw := w
	if h < w {
		sqw = h
	}
	//fmt.Println(sqw)
	//fmt.Println(sqw/2)
	//fmt.Println(w - sqw)
	for sqw > 2 {
		for x:=0;x<=w-sqw;x+=sqw-1 {
			for y:=0; y<=h-sqw;y+=sqw-1 {
				doSquareMidpoint(fullMap, x, y, sqw, 1, 2)

				//fmt.Printf("x: %d, y: %d, w: %d\n", x, y, sqw)
				//for i := 0; i < len(*fullMap); i++ {
				//	str := ""
				//	for j := 0; j < len((*fullMap)[0]); j++ {
				//		str += strconv.Itoa((*fullMap)[i][j])
				//	}
				//	fmt.Println(str)
				//}
				//fmt.Println("--------------")
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
}
