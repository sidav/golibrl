package RBR_generator

func (r *RBR) digCorridorIfPossible(x, y, dirx, diry, length int) bool {
	w := length * dirx
	h := length * diry
	wDec := 2 * dirx
	hDec := 2 * diry
	if w == 0 {
		wDec = 0
		w = 1
	}
	if h == 0 {
		hDec = 0
		h = 1
	}
	// TODO: allow corridors end in rooms or another corridors
	if r.isSpaceOfGivenType(x+dirx, y+diry, w-wDec, h-hDec, 1, TWALL) {
		// check if the end is not diagonally aligned to a floor
		corrEndX, corrEndY := x-dirx + length*dirx, y-diry + length*diry
		if r.countTiletypesAround(TFLOOR, corrEndX, corrEndY, false) > 0 ||
			r.countTiletypesAround(TFLOOR, corrEndX, corrEndY, true) == 0 {
			r.digSpace(x, y, w, h, 0)
			return true
		}
	}
	return false
}

func (r *RBR) placeCorridorFrom(x, y int) bool {
	// first, collect list of vectors of diggable directions near the x,y
	allowContinuation := rnd.RandInRange(0, 3) == 0
	dirs := r.pickListOfDiggableDirectionsFrom(x, y, allowContinuation)
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
		// try various lengths for each direction
		corrLength := rnd.RandInRange(r.MIN_CLENGTH, r.MAX_CLENGTH)
		for lenTry := 0; lenTry < r.MAX_CLENGTH; lenTry++ {
			digged = r.digCorridorIfPossible(x, y, vx, vy, corrLength)
			if digged || corrLength == r.MIN_CLENGTH {
				break
			}
			corrLength--
		}
		ind = (ind + 1) % len(*dirs)
		if ind == startind && !digged {
			return false
		}
	}
	r.placeDoorIfNeeded(x, y)
	return true
}
