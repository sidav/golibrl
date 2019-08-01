package main

import (
	"fmt"
	cw "github.com/sidav/golibrl/console"
	"github.com/sidav/golibrl/geometry"
)

func KeyToDirection(key string, diagsAllowed bool) (int, int) {
	var x, y int
	switch key {
	case "1", "b":
		x, y = -1, 1
	case "DOWN", "2", "j":
		x, y = 0, 1
	case "3", "n":
		x, y = 1, 1
	case "LEFT", "4", "h":
		x, y = -1, 0
	case "RIGHT", "6", "l":
		x, y = 1, 0
	case "7", "y":
		x, y = -1, -1
	case "UP", "8", "k":
		x, y = 0, -1
	case "9", "u":
		x, y = 1, -1
	}
	if !diagsAllowed && x*y != 0 {
		x = 0
		y = 0
	}
	return x, y
}


func angletest() {
	w, h := cw.GetConsoleSize()
	cx, cy := w/2, h/2
	vx, vy := 1, 0
	angle := 170
	key := ""
	for key != "ESCAPE" {
		cw.Clear_console()
		for i := 0; i < w; i++ {
			for j := 0; j < h; j++ {
				if geometry.AreCoordsInSector(i, j, cx, cy, vx, vy, angle) {
					cw.PutChar('+', i, j)
				}
			}
		}
		cw.PutString(fmt.Sprintf("angle %d, looking at %d, %d", angle, vx, vy), 0, 0)
		cw.Flush_console()
		key = cw.ReadKey()
		switch key {
		case ".": angle++
		case ",": angle--
		default:
			vx, vy = KeyToDirection(key, true)
		}
	}
}
