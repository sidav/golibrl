package main

import (
	"fmt"
	"github.com/sidav/golibrl/console"
	"github.com/sidav/golibrl/fov/basic_bresenham_fov"
	"github.com/sidav/golibrl/fov/basic_two_step_fov"
	"github.com/sidav/golibrl/fov/mill_fov"
	"github.com/sidav/golibrl/fov/permissive_fov"
	"github.com/sidav/golibrl/fov/strict_definition_fov"
	"github.com/sidav/golibrl/procedural_generation/CA_cave"
	"github.com/sidav/golibrl/procedural_generation/Fractal_landscape"
	"strconv"
	"strings"
	"time"
)

func main() {
	console.Init_console("test", console.TCellRenderer)
	defer console.Close_console()
	key := ""
	for key != "ESCAPE" {
		// angletest()
		// testFractalLandscape()
		// testCave()
		testFOV()
		console.Flush_console()
		key = console.ReadKey()
	}
}

var fovTestMap = &[]string{
	"#.#......##",
	"#.#....#..#",
	"#.#..#.....",
	"#.#...#..#.",
	"#.#....#...",
	"#.#........",
	"#..........",
	"#.#........",
	"#.#........",
	"#.#........",
	"#.#........",
	"#.#.....#..",
	"#.#........",
}

var cave *[]string

func opacityFunc(x, y int) bool {
	return rune((*cave)[x][y]) == '#'
}

func testFOV() {
	currentFovSelected := 1
	fovRadius := 15

	w, h := console.GetConsoleSize()
	cave = CA_cave.MakeCave(w, h, 40, 3, 25)
	//cave = fovTestMap
	//w, h = len(*fovTestMap), len((*fovTestMap)[0])

	px, py := w/2, h/2

	opacityMap := make([][]bool, w)

	for i := 0; i < len(*cave); i++ {
		opacityMap[i] = make([]bool, h)
		str := ' '
		for j := 0; j < len((*cave)[0]); j++ {
			str = rune((*cave)[i][j])
			if str == '#' {
				opacityMap[i][j] = true
			}
		}
	}

	key := ""
	for key != "ESCAPE" {
		// getVisibilityMap
		visMap, currentFovAlgorithmName := getVisMapAndNameForAlgorithm(currentFovSelected, px, py, fovRadius, w, h, &opacityMap)

		// render map
		for i := 0; i < len(*cave); i++ {
			for j := 0; j < len((*cave)[0]); j++ {
				str := rune((*cave)[i][j])
				if (*visMap)[i][j] {
					console.SetFgColor(console.WHITE)
					if str == '#' {
						console.SetFgColor(console.DARK_RED)
					}
					console.PutChar(str, i, j)
				} else {
					console.SetFgColor(console.BLACK)
					console.PutChar(' ', i, j)
				}
			}
		}
		console.SetFgColor(console.WHITE)
		console.PutChar('@', px, py)
		console.PutString(currentFovAlgorithmName, 0, 0)
		console.Flush_console()

		key = console.ReadKey()
		switch key {
		case "DOWN":
			py++
		case "UP":
			py--
		case "LEFT":
			px--
		case "RIGHT":
			px++
		case "-":
			fovRadius--
		case "+", "=":
			fovRadius++
		case "ENTER":
			// test-shmest
			fovAlgsPerfomanceCheck(px, py, w, h, fovRadius, &opacityMap)
		default:
			currentFovSelected, _ = strconv.Atoi(key)
		}
	}
}

func getVisMapAndNameForAlgorithm(currentFovSelected, px, py, fovRadius, w, h int, opacityMap *[][]bool) (*[][]bool, string) {
	var visMap *[][]bool
	currentFovAlgorithmName := "WTF FOV"
	switch currentFovSelected {
	case 1:
		visMap = strict_definition_fov.GetFovMapFrom(px, py, fovRadius, w, h, opacityFunc)
		currentFovAlgorithmName = "SDFOV"
	case 2:
		visMap = basic_bresenham_fov.GetFovMapFrom(px, py, fovRadius, w, h, opacityFunc)
		currentFovAlgorithmName = "Bresenham FOV"
	case 3:
		visMap = permissive_fov.GetFovMapFrom(px, py, fovRadius, w, h, opacityFunc)
		currentFovAlgorithmName = "Permissive FOV"
	case 4:
		visMap = mill_fov.GetFovMapFrom(px, py, fovRadius, w, h, opacityFunc)
		currentFovAlgorithmName = "Mill FOV"
	default:
		visMap = basic_two_step_fov.GetFovMapFrom(px, py, fovRadius, w, h, opacityFunc)
		currentFovAlgorithmName = "Two-step FOV"
	}
	return visMap, currentFovAlgorithmName
}

func fovAlgsPerfomanceCheck(px, py, w, h, fovRadius int, opacityMap *[][]bool) {
	const totalAlgs = 5
	const MillisecondsToTest = 1000
	console.SetBgColor(console.DARK_GRAY)
	console.SetFgColor(console.BLACK)
	borderW := 40
	border := strings.Repeat(" ", borderW)
	console.PutString(border, w/2 - borderW/2, h/2-totalAlgs-1)
	console.PutString(fmt.Sprintf("    Testing %d algs for %d ms:", totalAlgs, MillisecondsToTest), w/2 - borderW/2, h/2-totalAlgs-1)
	console.PutString(border, w/2 - borderW/2, h/2-totalAlgs)
	console.Flush_console()
	for i:=1; i<=totalAlgs;i++ {
		name:=""
		taken := 0
		start := time.Now()
		for time.Now().Sub(start) / time.Millisecond < MillisecondsToTest {
			_, name = getVisMapAndNameForAlgorithm(i, px, py, fovRadius, len(*opacityMap), len((*opacityMap)[0]), opacityMap)
			taken++
		}
		console.PutString(fmt.Sprintf(" %s:%d calculations; ", name, taken), w/2 - borderW/2, h/2-totalAlgs-1+i)
		console.PutString(border, w/2 - borderW/2, h/2-totalAlgs+i)
		console.Flush_console()
	}
	console.PutString(border, w/2 - borderW/2, h/2+1)
	console.PutString("    <Perfomance tests finished> ", w/2 - borderW/2, h/2+1)
	console.Flush_console()
	console.SetBgColor(console.BLACK)
	console.ReadKey()
}

func testCave() {
	w, h := console.GetConsoleSize()
	cave := CA_cave.MakeCave(w, h, 40, 4, -1)
	console.SetFgColor(console.WHITE)
	for i := 0; i < len(*cave); i++ {
		str := ' '
		for j := 0; j < len((*cave)[0]); j++ {
			str = rune((*cave)[i][j])
			switch str {
			case '#':
				console.SetFgColor(console.DARK_RED)
			default:
				console.SetFgColor(console.WHITE)
			}
			console.PutChar(str, i, j)
		}
	}
}

func testFractalLandscape() {
	land := Fractal_landscape.GenHeightMap(129, 65)
	// return

	console.SetFgColor(console.WHITE)
	for i := 0; i < len(*land); i++ {
		str := ' '
		for j := 0; j < len((*land)[0]); j++ {
			switch cur := (*land)[i][j]; {
			case cur < -10:
				str = '~'
				console.SetFgColor(console.DARK_BLUE)
			case cur < 0:
				str = '~'
				console.SetFgColor(console.BLUE)
			case cur < 9:
				str = '.'
				console.SetFgColor(console.YELLOW)
			case cur < 22:
				str = ','
				console.SetFgColor(console.DARK_YELLOW)
			case cur < 40:
				str = 'T'
				console.SetFgColor(console.GREEN)
			case cur < 50:
				str = '^'
				console.SetFgColor(console.DARK_GRAY)
			default:
				str = '^'
				console.SetFgColor(console.WHITE)
			}
			console.PutChar(str, i, j)
		}
	}
}
