package RBR_generator

func (r *RBR) digSpace(x, y, w, h int) {
	if w < 0 {
		x = x + w + 1
		w = -w
	}
	if h < 0 {
		y = y + h + 1
		h = -h
	}
	for cx := x; cx < x+w; cx++ {
		for cy := y; cy < y+h; cy++ {
			r.tiles[cx][cy].tiletype = TFLOOR
		}
	}
}

func (r *RBR) countTiletypesAround(ttype byte, x, y int, diagonals bool) int {
	ttypes := 0
	for vx := -1; vx <= 1; vx++ {
		for vy := -1; vy <= 1; vy++ {
			if vx == vy && vx == 0 {
				continue
			}
			if !diagonals && vx*vy != 0 {
				continue
			}
			cx := x + vx
			cy := y + vy
			if cx < 0 || cy < 0 || cx >= r.mapw || cy >= r.maph {
				continue
			}
			if r.tiles[cx][cy].tiletype == ttype {
				ttypes++
			}
		}
	}
	return ttypes
}
func (r *RBR) isSpaceOfGivenType(x, y, w, h, outlineThickness int, ttype byte) bool {
	if w < 0 {
		x = x + w + 1
		w = -w
	}
	if h < 0 {
		y = y + h + 1
		h = -h
	}
	x -= outlineThickness
	y -= outlineThickness
	w += 2 * outlineThickness
	h += 2 * outlineThickness
	if x < 0 || y < 0 || x+w >= r.mapw || y+h >= r.maph {
		return false
	}
	for cx := x; cx < x+w; cx++ {
		for cy := y; cy < y+h; cy++ {
			if r.tiles[cx][cy].tiletype != ttype {
				return false
			}
		}
	}
	return true
}

func (r *RBR) digRoomIfPossible(x, y, w, h, oulinethick int) bool {
	if r.isSpaceOfGivenType(x, y, w, h, 1, TWALL) {
		r.digSpace(x, y, w, h)
		return true
	}
	return false
}
