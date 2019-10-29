package RBR_generator

func (r *RBR) tryPlaceRoom(x, y, dirx, diry, roomw, roomh, roomId int) bool {
	w := roomw * dirx
	h := roomh * diry
	if w == 0 {
		w = roomw
	}
	if h == 0 {
		h = roomh
	}
	if r.isSpaceOfGivenType(x+dirx, y+diry, w, h, 1, TWALL) {
		r.digSpace(x+dirx, y+diry, w, h, roomId)
		return true
	}
	return false
}

func (r *RBR) placeRoomFromJunction(x, y, roomId int) bool {
	// first, collect list of vectors of diggable directions near the x,y
	dirs := r.pickListOfDiggableDirectionsFrom(x, y, true)
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
		digged = r.tryPlaceRoom(x, y, vx, vy, roomW, roomH, roomId)
		ind = (ind + 1) % len(*dirs)
		if ind == startind && !digged {
			return false
		}
	}
	r.tiles[x][y].tiletype = TDOOR
	return true
}

// Experimental 

func (r *RBR) placeRoomByPicking(roomId int, deadendOnly bool) bool {
	placeFound := false
	tries := 0 
	maxtries := r.MAX_RSIZE * r.MAX_RSIZE - r.MIN_RSIZE*r.MIN_RSIZE
finding_place:
	for tries < maxtries {
		tries++
		roomW, roomH := rnd.RandInRange(r.MIN_RSIZE, r.MAX_RSIZE), rnd.RandInRange(r.MIN_RSIZE, r.MAX_RSIZE)
		coordsList := r.pickListOfCoordinatesForRoomToBeFit(roomW, roomH)
		if coordsList == nil {
			continue finding_place
		}
		selectedCoordIndex := rnd.Rand(len(*coordsList))
		currCoordIndex := selectedCoordIndex
	trying_coords:
		for {
			x, y := (*coordsList)[currCoordIndex][0], (*coordsList)[currCoordIndex][1]
			jx, jy := r.pickJunctionTileForPotentialRoom(x, y, roomW, roomH, deadendOnly)
			if jx != -1 && jy != -1 {
				r.digSpace(x, y, roomW, roomH, roomId)
				r.tiles[jx][jy].tiletype = TDOOR
				placeFound = true 
				break finding_place
			}
			currCoordIndex = (currCoordIndex+1) % len(*coordsList)
			if currCoordIndex == selectedCoordIndex {
				break trying_coords
			}
		}
	}
	return placeFound
}
