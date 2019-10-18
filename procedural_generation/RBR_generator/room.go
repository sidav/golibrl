package RBR_generator

func (r *RBR) tryPlaceRoom(x, y, dirx, diry, roomw, roomh int) bool {
	w := roomw * dirx
	h := roomh * diry
	if w == 0 {
		w = roomw
	}
	if h == 0 {
		h = roomh
	}
	if r.isSpaceOfGivenType(x+dirx, y+diry, w, h, 1, TWALL) {
		r.digSpace(x+dirx, y+diry, w, h)
		return true
	}
	return false
}

func (r *RBR) placeRoomFromJunction(x, y int) bool {
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
			roomW := rnd.RandInRange(r.MIN_RSIZE, r.MAX_RSIZE)
			roomH := rnd.RandInRange(r.MIN_RSIZE, r.MAX_RSIZE)
			digged = r.tryPlaceRoom(x, y, vx, vy, roomW, roomH)
			ind = (ind + 1) % len(*dirs)
			if ind == startind && !digged {
				return false
			}
		}
		r.tiles[x][y].tiletype = TDOOR
		return true
}
