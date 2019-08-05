package basic_two_step_fov

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

func GetFovMapFrom(fromx, fromy, radius int) *[][]bool {
	emptyVisibilityMap(len(*opaque), len((*opaque)[0]))
	doFirstStep(fromx, fromy, radius)
	doSecondStep(fromx, fromy, radius)
	return visible
}

func doFirstStep(fromx, fromy, radius int) {
	circle := graphic_primitives.GetApproxCircleAroundRect(fromx, fromy, 0, 0, radius)
	for i := range *circle {
		line := graphic_primitives.GetLine(fromx, fromy, (*circle)[i].X, (*circle)[i].Y)
		for j := range *line {
			lx, ly := (*line)[j].X, (*line)[j].Y
			if geometry.AreCoordsInRect(lx, ly, 0, 0, len(*visible), len((*visible)[0])) {
				(*visible)[lx][ly] = true
				if j > 0 && (*opaque)[lx][ly] {
					break
				}
			}
		}
	}
}

func doSecondStep(fromx, fromy, radius int) {
	visibleList := make([]graphic_primitives.Point, 0)
	for x := fromx - radius + 1; x < fromx+radius-1; x++ {
		for y := fromy - radius + 1; y < fromy+radius-1; y++ {
			if geometry.AreCoordsInRect(x, y, 1, 1, len(*opaque)-2, len((*opaque)[0])-2) {
				totalVisibles := 0
			checkNeighbours:
				for i := -1; i <= 1; i++ {
					for j := -1; j <= 1; j++ {
						if i*j == 0 && i != j && (*visible)[x+i][y+j] {
							totalVisibles++
						}
						if totalVisibles > 2 {
							visibleList = append(visibleList, graphic_primitives.Point{x, y})
							break checkNeighbours
						}
					}
				}
			}
		}
	}
	for i := range visibleList {
		x, y := visibleList[i].GetCoords()
		(*visible)[x][y] = true
	}
}
