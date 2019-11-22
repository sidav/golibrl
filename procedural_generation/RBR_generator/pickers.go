package RBR_generator

func (r *RBR) pickListOfDiggableDirectionsFrom(x, y int, allowContinuation bool) *[][]int {
	dirs := make([][]int, 0)
	for vx := -1; vx <= 1; vx++ {
		for vy := -1; vy <= 1; vy++ {
			if x+vx < 0 || y+vy < 0 || x+vx >= r.mapw || y+vy >= r.maph {
				continue
			}
			if vx != vy && vx*vy == 0 && r.tiles[x+vx][y+vy].TileType == TWALL {
				if allowContinuation || r.tiles[x-vx][y-vy].TileType != TFLOOR {
					dirs = append(dirs, []int{vx, vy})
				}
			}
		}
	}
	return &dirs
}

func (r *RBR) isTileSuitableForJunction(x, y int, deadendOnly bool) bool {
	if r.tiles[x][y].TileType == TWALL {
		if deadendOnly {
			walls := r.countTiletypesAround(TWALL, x, y, true)
			floors := r.countTiletypesAround(TFLOOR, x, y, false)
			if walls == 7 && floors == 1 {
				return true
			}
		} else {
			walls := r.countTiletypesAround(TWALL, x, y, false)
			floors := r.countTiletypesAround(TFLOOR, x, y, false)
			if walls == 3 && floors == 1 {
				return true
			}
		}
	}
	return false
}

func (r *RBR) pickJunctionTile(fromx, fromy, tox, toy int, deadendOnly bool) (int, int) {
	listOfAppropriateCoords := make([][]int, 0)
	for x := fromx; x < tox; x++ {
		for y := fromy; y < toy; y++ {
			if r.isTileSuitableForJunction(x, y, deadendOnly) {
				listOfAppropriateCoords = append(listOfAppropriateCoords, []int{x, y})
			}
		}
	}
	if len(listOfAppropriateCoords) == 0 {
		return -1, -1
	}
	indx := rnd.Rand(len(listOfAppropriateCoords))
	return listOfAppropriateCoords[indx][0], listOfAppropriateCoords[indx][1]
}

// Experimental:

func (r *RBR) pickJunctionTileForPotentialRoom(rx, ry, w, h int, deadendOnly bool) (int, int) {
	listOfAppropriateCoords := make([][]int, 0)
	for x := rx; x < rx+w; x++ {
		if r.isTileSuitableForJunction(x, ry-1, deadendOnly) {
			listOfAppropriateCoords = append(listOfAppropriateCoords, []int{x,ry-1})
		}
		if r.isTileSuitableForJunction(x, ry+h, deadendOnly) {
			listOfAppropriateCoords = append(listOfAppropriateCoords, []int{x,ry+h})
		}
	}
	for y := ry; y < ry+h; y++ {
		if r.isTileSuitableForJunction(rx-1, y, deadendOnly) {
			listOfAppropriateCoords = append(listOfAppropriateCoords, []int{rx-1, y})
		}
		if r.isTileSuitableForJunction(rx+w, y, deadendOnly) {
			listOfAppropriateCoords = append(listOfAppropriateCoords, []int{rx+w, y})
		}
	}

	if len(listOfAppropriateCoords) == 0 {
		return -1, -1 
	}
	coord := listOfAppropriateCoords[rnd.Rand(len(listOfAppropriateCoords))] 
	return coord[0], coord[1]
}

func (r *RBR) pickListOfCoordinatesForRoomToBeFit(w, h int) *[][]int {
	listOfPotentiallyAppropriateCoords := make([][]int, 0)
	for x := 2; x+w<r.mapw-1; x++ {
		for y := 2; y+h < r.maph-1; y++ {
			if r.isSpaceOfGivenType(x, y, w, h, 1, TWALL) {
				listOfPotentiallyAppropriateCoords = append(listOfPotentiallyAppropriateCoords, []int{x, y})
			}
		}
	}
	if len(listOfPotentiallyAppropriateCoords) == 0 {
		return nil 
	}
	return &listOfPotentiallyAppropriateCoords
}


