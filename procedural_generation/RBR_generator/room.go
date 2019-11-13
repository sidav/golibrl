package RBR_generator

// Experimental 

func (r *RBR) placeRoomByPicking(roomId int, deadendOnly, placeVault bool) bool {
	placeFound := false
	tries := 0 
	maxtries := 1 // Not needed? // r.MAX_RSIZE * r.MAX_RSIZE - r.MIN_RSIZE*r.MIN_RSIZE
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
				r.tryPlaceVaultOfGivenSizeAtCoords(x+1, y+1, roomW-2, roomH-2)
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
