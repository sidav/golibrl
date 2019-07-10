package strict_definition_fov

// IMPLEMENTED IN GOLANG FROM http://www.roguebasin.com/index.php?title=LOS_using_strict_definition
// Optimized by sidav

import (
	"github.com/sidav/golibrl/geometry"
	"github.com/sidav/golibrl/graphic_primitives"
	"math"
)

var (
	mapw, maph int
	opaque  *[][]bool
	visible *[][]bool
)

func SetOpacityMap(o *[][]bool) {
	opaque = o
	mapw = len(*opaque)
	maph = len((*opaque)[0])
}

func emptyVisibilityMap(w, h int) {
	vis := make([][]bool, w)
	for i := range vis {
		vis[i] = make([]bool, h)
	}
	visible = &vis
}

func Fov(x, y, radius int) *[][]bool {
	radius++
	emptyVisibilityMap(len(*opaque), len((*opaque)[0]))
	circle := graphic_primitives.GetApproxCircleAroundRect(x, y, 0, 0, radius)
	for i := range *circle {
		los(x, y, (*circle)[i].X, (*circle)[i].Y)
	}
	return visible
}

/* Los calculation */
func los(x0, y0, x1, y1 int) {
	var sx, sy, xnext, ynext, dx, dy int
	var dist float64

	dx = x1 - x0
	dy = y1 - y0

	//determine which quadrant to we're calculating: we climb in these two directions
	sx = -1
	sy = -1
	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}
	xnext = x0
	ynext = y0

	//calculate length of line to cast (distance from start to final tile)
	dist = math.Sqrt(float64(dx*dx + dy*dy))
	for xnext != x1 || ynext != y1 { //essentially casting a ray of length radius: (radius^3)

		if geometry.AreCoordsInRect(xnext, ynext, 0, 0, mapw, maph) {
			(*visible)[xnext][ynext] = true
			if (*opaque)[xnext][ynext] { // or any equivalent
				// tag_memorised(xnext, ynext); // make a note of the wall
				return
			}
		}

		// Line-to-point distance formula < 0.5
		if math.Abs(float64(dy*(xnext-x0+sx)-dx*(ynext-y0)))/dist < 0.5 {
			xnext += sx
		} else if math.Abs(float64(dy*(xnext-x0)-dx*(ynext-y0+sy)))/dist < 0.5 {
			ynext += sy
		} else {
			xnext += sx
			ynext += sy
		}
	}
	if geometry.AreCoordsInRect(x1, y1, 0, 0, mapw, maph) {
		(*visible)[x1][y1] = true
	}
}
