package RBR_generator

func (r *RBR) digCorridorIfPossible(x, y, dirx, diry, length int) bool {
	w := length * dirx
	h := length * diry
	if w == 0 {
		w = 1
	}
	if h == 0 {
		h = 1
	}
	// TODO: allow corridors end in rooms or another corridors 
	if r.isSpaceOfGivenType(x+dirx, y+diry, w, h, 1, TWALL) {
		r.digSpace(x, y, w, h)
		return true
	}
	return false
}

func (r *RBR) pickTileForCorridorPlacement() (int, int) {
	listOfAppropriateCoords := make([][]int, 0)
	for x := 0; x < r.mapw; x++ {
		for y := 0; y < r.maph; y++ {
			walls := r.countTiletypesAround(TWALL, x, y, false)
			floors := r.countTiletypesAround(TFLOOR, x, y, false)
			if r.tiles[x][y].tiletype == TWALL && walls == 3 && floors == 1 {
				listOfAppropriateCoords = append(listOfAppropriateCoords, []int{x, y})
			}
		}
	}
	if len(listOfAppropriateCoords) == 0 {
		panic("Oh fuck.")
	}
	indx := rnd.Rand(len(listOfAppropriateCoords))
	return listOfAppropriateCoords[indx][0], listOfAppropriateCoords[indx][1]
}

func (r *RBR) placeCorridorFrom(x, y int) bool {
	// first, collect list of vectors of diggable directions near the x,y
	dirs := make([][]int, 0)
	for vx := -1; vx <= 1; vx++ {
		for vy := -1; vy <= 1; vy++ {
			if x+vx < 0 || y+vy < 0 || x+vx >= r.mapw || y+vy >= r.maph {
				continue
			}
			if vx != vy && vx*vy == 0 && r.tiles[x+vx][y+vy].tiletype == TWALL {
				dirs = append(dirs, []int{vx, vy})
			}
		}
	}
	if len(dirs) == 0 {
		return false
	}
	// next, let's pick a random vector from them
	startind := rnd.Rand(len(dirs))
	// ...starting from that index, try every direction.
	ind := startind
	digged := false
	for !digged {
		vx, vy := dirs[ind][0], dirs[ind][1]
		corrLength := rnd.RandInRange(r.MIN_CLENGTH, r.MAX_CLENGTH)
		digged = r.digCorridorIfPossible(x, y, vx, vy, corrLength)
		ind = (ind + 1) % len(dirs)
		if ind == startind && !digged {
			return false
		}
	}
	r.tiles[x][y].tiletype = TDOOR
	return true
}
