package RBR_generator

func (r *RBR) digSpace(x, y, w, h int) {
	for cx := x; cx < x+w; cx++ {
		for cy := y; cy < y+h; cy++ {
			r.tiles[cx][cy].tiletype = TFLOOR
		}
	}
}

func (r *RBR) isSpaceOfGivenType(x, y, w, h int, ttype byte) bool {
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
	if r.isSpaceOfGivenType(x-oulinethick, y-oulinethick, w+2*oulinethick, h+2*oulinethick, TWALL) {
		r.digSpace(x, y, w, h)
		return true
	}
	return false
}

func (r *RBR) digCorridorIfPossible(x, y, dirx, diry) bool {
	if r.isSpaceOfGivenType(x-oulinethick, y-oulinethick, w+2*oulinethick, h+2*oulinethick, TWALL) {
		r.digSpace(x, y, w, h)
		return true
	}
	return false
}

func (r *RBR) digCorridorNear(x, y int) {
	// first, collect list of vectors of diggable directions near the x,y 
	dirs := make([][]int, 0)
	for vx := -1; vx <= 1; vx++ {
		for vy := -1; vy <= 1; vy++ {
			if vx != vy && vx * vy == 0 && r.tiles[x+vx][y+vy].tiletype == TWALL {
				dirs = append(dirs, []int{vx, vy})
			}
		}
	}
	if len(dirs) == 0 {
		return 
	}
	// next, let's pick a random vector from them 
	ind := rnd.RandInRange(0, len(dirs))
	vx, vy := dirs[ind][0], dirs[ind][1]
	// ...aaaand use it. 

}
