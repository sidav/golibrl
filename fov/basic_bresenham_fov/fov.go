package basic_bresenham_fov

import (
	"github.com/sidav/golibrl/geometry"
	"github.com/sidav/golibrl/graphic_primitives"
)

var (
	opaque  *[][]bool
	visible *[][]bool
)

func SetOpacityMap(o *[][]bool) {
	opaque = o
}

func emptyVisibilityMap(w, h int) {
	vis := make([][]bool, w)
	for i := range vis {
		vis[i] = make([]bool, h)
	}
	visible = &vis
}

func GetCircleVisibilityMap(fromx, fromy, radius int) *[][]bool {
	emptyVisibilityMap(len(*opaque), len((*opaque)[0]))
	doFirstStep(fromx, fromy, radius)
	// doSecondStep(fromx, fromy, radius)
	return visible
}

func doFirstStep(fromx, fromy, radius int) {
	for i:=fromx-radius; i < fromx+radius;i++ {
		for j:=fromy-radius; j < fromy+radius;j++ {
			line := graphic_primitives.GetLine(fromx, fromy, i, j)
			for lineIndex := range *line {
				lx, ly := (*line)[lineIndex].X, (*line)[lineIndex].Y
				if geometry.AreCoordsInRect(lx, ly, 0, 0, len(*visible), len((*visible)[0])) {
					(*visible)[lx][ly] = true
					if lineIndex > 0 && (*opaque)[lx][ly] {
						break
					}
				}
			}
		}
	}
}
