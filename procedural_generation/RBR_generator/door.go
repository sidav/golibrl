package RBR_generator

func (r *RBR) placeDoorIfNeeded(x, y int) {
	if r.isTileAdjacentToDifferentRoomIDs(x, y) {
		r.tiles[x][y].tiletype = TDOOR 
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
