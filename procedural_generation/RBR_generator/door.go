package RBR_generator

func (r *RBR) placeDoorIfNeeded(x, y int) {
	if r.isTileAdjacentToDifferentRoomIDs(x, y) || r.isTileAdjacentToDifferentSecAreas(x, y) {
		r.tiles[x][y].TileType = TDOOR
	}
}

func (r *RBR) finalizeDoorsSecArea() {
	for x := 0; x < r.mapw; x++ {
		for y := 0; y < r.maph; y++ {
			if r.tiles[x][y].TileType == TDOOR {
				if r.isTileAdjacentToDifferentSecAreas(x, y) {
					r.tiles[x][y].SecArea = r.getHighestSecAreaNearTile(x, y)
				} else {
					r.tiles[x][y].SecArea = 0 // don't lock doors which aren't connecting two sec areas
				}
			}
		}
	}
}

func (r *RBR) placeRandomDoors(doorsNum int) {
	for door := 0; door < doorsNum; door++ {
		suitableDoorCoords := make([][]int, 0)
		for x := 1; x < r.mapw; x++ {
			for y := 1; y < r.maph; y++ {
				if r.tiles[x][y].TileType == TWALL && r.countTiletypesAround(TFLOOR, x, y, false) == 2 && r.countTiletypesAround(TDOOR, x, y, false) == 0 {
					suitableDoorCoords = append(suitableDoorCoords, []int{x, y})
				}
			}
		}
		if len(suitableDoorCoords) == 0 {
			return
		}
		coords := suitableDoorCoords[rnd.Rand(len(suitableDoorCoords))]
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
			if vx != vy && vx*vy == 0 && r.tiles[x+vx][y+vy].TileType != TWALL {
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

func (r *RBR) getHighestSecAreaNearTile(x, y int) int16 {
	var currSecArea int16 = 0
	for vx := -1; vx <= 1; vx++ {
		for vy := -1; vy <= 1; vy++ {
			if x+vx < 0 || y+vy < 0 || x+vx >= r.mapw || y+vy >= r.maph {
				continue
			}
			if vx != vy && vx*vy == 0 && r.tiles[x+vx][y+vy].SecArea > currSecArea {
				currSecArea++
			}
		}
	}
	return currSecArea
}

func (r *RBR) isTileAdjacentToDifferentSecAreas(x, y int) bool {
	var currSecArea int16 = -1
	for vx := -1; vx <= 1; vx++ {
		for vy := -1; vy <= 1; vy++ {
			if x+vx < 0 || y+vy < 0 || x+vx >= r.mapw || y+vy >= r.maph {
				continue
			}
			if vx != vy && vx*vy == 0 && r.tiles[x+vx][y+vy].TileType != TWALL {
				if currSecArea == -1 {
					currSecArea = r.tiles[x+vx][y+vy].SecArea
					continue
				}
				if r.tiles[x+vx][y+vy].SecArea != currSecArea {
					return true
				}
			}
		}
	}
	return false
}
