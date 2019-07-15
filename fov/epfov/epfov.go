package epfov

var (
	mapw, maph, rangeLimit int
	opaque                 *[][]bool
	visible                *[][]bool
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

///

const (
	STEP_SIZE = 16
)

var (
	offset, limit int
	current_view  *view_t
	views         *view_t
	bumpidx       = 0
)

func is_blocked(view *view_t, startX, startY, x, y, dx, dy int, light_walls bool) bool {
	posx := x*dx/STEP_SIZE + startX
	posy := y*dy/STEP_SIZE + startY
	blocked := (*opaque)[posx][posy]
	if !blocked || light_walls {
		(*visible)[posx][posy] = true
	}
	return blocked
}
