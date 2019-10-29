package RBR_generator

func (r *RBR) pickListOfDiggableDirectionsFrom(x, y int, allowContinuation bool) *[][]int {
	dirs := make([][]int, 0)
	for vx := -1; vx <= 1; vx++ {
		for vy := -1; vy <= 1; vy++ {
			if x+vx < 0 || y+vy < 0 || x+vx >= r.mapw || y+vy >= r.maph {
				continue
			}
			if vx != vy && vx*vy == 0 && r.tiles[x+vx][y+vy].tiletype == TWALL {
				if allowContinuation || r.tiles[x-vx][y-vy].tiletype != TFLOOR {
					dirs = append(dirs, []int{vx, vy})
				}
			}
		}
	}
	return &dirs
}

func (r *RBR) isTileSuitableForJunction(x, y int, deadendOnly bool) bool {
	if r.tiles[x][y].tiletype == TWALL {
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

// func (r *RBR) pickJunctionTileForRoom(rx, ry, w, h int) (int, int) {
// 	listOfAppropriateCoords := make([][]int, 0)
// 	for x := rx - 1; x <= rx+w; x++ {
// 		for y := ry - 1; y <= ry+h; y++ {
// 			if (x == rx-1 || x == rx+w) && (y == ry-1 || y == ry+h) {
// 				if r.isTileSuitableForJunction(x, y, false) {
// 					listOfAppropriateCoords = append(listOfAppropriateCoords, []int{x,y})
// 				}
// 			}
// 		}
// 	}
// 	if len(listOfAppropriateCoords) == 0 {
// 		return -1, -1 
// 	}
// 	coord := listOfAppropriateCoords[rnd.Rand(len(listOfAppropriateCoords))] 
// 	return coord[0], coord[1]
// }
