package RBR_generator

func (r *RBR) placeDoorIfNeeded(x, y int) {
	if r.isTileAdjacentToDifferentRoomIDs(x, y) {
		r.tiles[x][y].tiletype = TDOOR 
	}
}

func ( r* RBR) placeRandomDoors(doorsNum int) {
	for door := 0; door < doorsNum; door++ {
		suitableDoorCoords := make([][]int, 0)
		for x := 1; x < r.mapw; x++ {
			for y := 1; y < r.maph; y++ {
				if r.tiles[x][y].tiletype == TWALL && r.countTiletypesAround(TFLOOR, x, y, false) == 2 && r.countTiletypesAround(TDOOR, x, y, false) == 0 {
					suitableDoorCoords = append(suitableDoorCoords, []int{x, y})
				}
			}
		}
		if len(suitableDoorCoords) == 0 {
			return 
		}
		coords := suitableDoorCoords[rnd.Rand(len(suitableDoorCoords))]
		r.tiles[0][0].tiletype = TDOOR
		r.placeDoorIfNeeded(coords[0], coords[1])
	}
}

func (r *RBR) isTileAdjacentToDifferentRoomIDs(x, y int) bool {
	currId := -1 
	for vx := -1; vx <= 1; vx++ {
		for vy := -1; vy <= 1; vy++ {
			if x+vx < 0 || y+vy < 0 || x+vx >= r.mapw || y+vy >= r.maph {
				continue
			}
			if vx != vy && vx*vy == 0 {
				if currId == -1 {
					currId = r.tiles[x+vx][y+vy].roomId
					continue 
				}
				if r.tiles[x+vx][y+vy].roomId != currId {
					return true 
				}
			}
		}
	}
	return false
}
