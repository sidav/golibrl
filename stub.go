package main

import (
	"github.com/sidav/golibrl/console"
	"github.com/sidav/golibrl/fov/basic_two_step_fov"
	"github.com/sidav/golibrl/fov/strict_definition_fov"
	"github.com/sidav/golibrl/procedural_generation/CA_cave"
	"github.com/sidav/golibrl/procedural_generation/Fractal_landscape"
	"strconv"
)

func main() {
	console.Init_console("test", console.TCellRenderer)
	defer console.Close_console()
	key := ""
	for key != "ESCAPE" {
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

func testFOV() {
	currentFovSelected := 1
	currentFovAlgorithmName := "UNSELECTED"
	fovRadius := 15

	w, h := console.GetConsoleSize()
	cave := CA_cave.MakeCave(w, h, 3, 1)
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
		var visMap *[][]bool
		switch currentFovSelected {
		case 1:
			strict_definition_fov.SetOpacityMap(&opacityMap)
			visMap = strict_definition_fov.Fov(px, py, fovRadius)
			currentFovAlgorithmName = "Strict Definition FOV"
		default:
			basic_two_step_fov.SetOpacityMap(&opacityMap)
			visMap = basic_two_step_fov.GetCircleVisibilityMap(px, py, fovRadius)
			currentFovAlgorithmName = "Two-step FOV"
		}

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
		case "1", "2", "3":
			currentFovSelected, _ = strconv.Atoi(key)
		}
	}
}

func testCave() {
	cave := CA_cave.MakeCave(60, 20, 3, -1)
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
