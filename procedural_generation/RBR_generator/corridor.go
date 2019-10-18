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

func (r *RBR) placeCorridorFrom(x, y int) bool {
	// first, collect list of vectors of diggable directions near the x,y
	dirs := r.pickListOfDiggableDirectionsFrom(x, y)
	if len(*dirs) == 0 {
		return false
	}
	// next, let's pick a random vector from them
	startind := rnd.Rand(len(*dirs))
	// ...starting from that index, try every direction.
	ind := startind
	digged := false
	for !digged {
		vx, vy := (*dirs)[ind][0], (*dirs)[ind][1]
		corrLength := rnd.RandInRange(r.MIN_CLENGTH, r.MAX_CLENGTH)
		digged = r.digCorridorIfPossible(x, y, vx, vy, corrLength)
		ind = (ind + 1) % len(*dirs)
		if ind == startind && !digged {
			return false
		}
	}
	r.tiles[x][y].tiletype = TDOOR
	return true
}
