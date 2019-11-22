package RBR_generator

func (r *RBR) placeStairs(numUp, numDown int, maximizeOutEntranceSecLevel bool) {
	for curr := 0; curr < numUp; curr++ {
		placed := r.placeStairAtRandom(TPREVLEVELSTAIR, 0, 0)
		if !placed {
			panic("RBR_GEN: no entrance placed")
		}
	}
	minsecid := int16(0) 
	maxsecid := int16(r.NUM_SEC_AREAS-1)
	if maximizeOutEntranceSecLevel {
		minsecid = maxsecid
	}
	for curr := 0; curr < numDown; curr++ {
		r.placeStairAtRandom(TNEXTLEVELSTAIR, minsecid, maxsecid)
	}
}

func (r *RBR) placeStairAtRandom(ttype byte, minSecArea, maxSecArea int16) bool {
	suitableEntrCoords := make([][]int, 0)
	for x := 1; x < r.mapw; x++ {
		for y := 1; y < r.maph; y++ {
			if r.tiles[x][y].TileType == TFLOOR && 
			r.tiles[x][y].SecArea >= minSecArea &&
			r.tiles[x][y].SecArea <= maxSecArea &&
			 r.countTiletypesAround(TFLOOR, x, y, true) > 2 {
				suitableEntrCoords = append(suitableEntrCoords, []int{x, y})
			}
		}
	}
	if len(suitableEntrCoords) == 0 {
		return false 
	}
	coords := suitableEntrCoords[rnd.Rand(len(suitableEntrCoords))]
	r.tiles[coords[0]][coords[1]].setProperties(ttype, -1, -1)
	return true 
}
