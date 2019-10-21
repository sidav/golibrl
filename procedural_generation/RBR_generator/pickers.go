package RBR_generator

func (r *RBR) pickListOfDiggableDirectionsFrom(x, y int) *[][]int {
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
	return &dirs 
}

func (r *RBR) pickJunctionTile(fromx, fromy, tox, toy int, deadendOnly bool) (int, int) {
	listOfAppropriateCoords := make([][]int, 0)
	for x := fromx; x < tox; x++ {
		for y := fromy; y < toy; y++ {
			if deadendOnly {
				walls := r.countTiletypesAround(TWALL, x, y, true)
				floors := r.countTiletypesAround(TFLOOR, x, y, false)
				if r.tiles[x][y].tiletype == TWALL && walls == 7 && floors == 1 {
					listOfAppropriateCoords = append(listOfAppropriateCoords, []int{x, y})
				}
			} else {
				walls := r.countTiletypesAround(TWALL, x, y, false)
				floors := r.countTiletypesAround(TFLOOR, x, y, false)
				if r.tiles[x][y].tiletype == TWALL && walls == 3 && floors == 1 {
					listOfAppropriateCoords = append(listOfAppropriateCoords, []int{x, y})
				}
			}
		}
	}
	if len(listOfAppropriateCoords) == 0 {
		return -1, -1
	}
	indx := rnd.Rand(len(listOfAppropriateCoords))
	return listOfAppropriateCoords[indx][0], listOfAppropriateCoords[indx][1]
}
