package procedural_generation

import (
	"time"
)

// cellular automata based cave generator.

var seed int

const (
	a = 2416
	c = 374441
	// m = 1771875
	m = 1000000005721 // prime
)

var (
	x int
)

func randomize() int {
	x = int(time.Duration(time.Now().UnixNano())/time.Millisecond) % m
	return x
}

func random(modulo int) int {
	x = (x*a + c) % m
	if modulo != 0 {
		return x % modulo
	} else {
		return x
	}
}

func MakeCave(w, h, smoothness, seed int) *[]string {
	x = seed
	if seed < 0 {
		randomize()
	}
	if smoothness < 1 {
		smoothness = 1
	}
	cave := randomInitialFill(w, h)
	cave = cycle(cave, smoothness)
	return cave
}

func randomInitialFill(w, h int) *[]string {
	const WALL_PERCENTAGE = 35
	cave := make([]string, h)
	for ind := range cave {
		for i := 0; i < w; i++ {
			if ind - h/2 <= 2 && ind - h/2 > -2 {
				cave[ind] += "."
				continue
			}
			rnd := random(100)
			if rnd < WALL_PERCENTAGE {
				cave[ind] += "#"
			} else {
				cave[ind] += "."
			}
		}
	}
	return &cave
}

func countWallsInRange(x, y, r int, cave *[]string) int { // diagonals are equidistant to straight dirs
	total := 0
	for i := x - r; i <= x+r; i++ {
		for j := y - r; j <= y+r; j++ {
			if i < 0 || j < 0 || i >= len(*cave) || j >= len((*cave)[0]) {
				total++ // out of bounds counts as a wall!
				continue
			}
			// fmt.Printf("ij %d,%d while %d, %d", i, j, len(*cave), len((*cave)[0]))
			if (*cave)[i][j] == '#' {
				total++
			}
		}
	}
	return total
}

func cycle(inpCave *[]string, max int) *[]string {
	cave := *inpCave
	w := len(cave)
	h := len((cave)[0])
	var newCave []string
	changesMade := true
	loopNum := 0
	for changesMade && loopNum < max {

		//for _, s := range cave {
		//	fmt.Println(s)
		//}
		//fmt.Println("-----------------")

		changesMade = false
		newCave = make([]string, w)

		for ind := range newCave {
			for i := 0; i < h; i++ {
				sym := "."
				walls1 := countWallsInRange(ind, i, 1, &cave)
				walls2 := countWallsInRange(ind, i, 2, &cave)
				if walls1 >= 5 || walls2 <= 1 {
					sym = "#"
				}
				if (cave)[ind][i] != sym[0] {
					changesMade = true
				}
				newCave[ind] += sym
			}
		}
		cave = newCave
		loopNum++
	}
	return &cave
}
