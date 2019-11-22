package RBR_generator

func (r *RBR) pickListOfCoordinatesForGivenVaultToBeFit(vaultstrs *[]string) *[][]int {
	h, w := len(*vaultstrs), len((*vaultstrs)[0])
	listOfPotentiallyAppropriateCoords := make([][]int, 0)
	for x := 2; x+w < r.mapw-1; x++ {
global:
		for y := 2; y+h < r.maph-1; y++ {
			for vx := 0; vx < w; vx++ {
				for vy := 0; vy < h; vy++ {
					if !(r.tiles[x+vx][y+vy].TileType == TWALL || rune((*vaultstrs)[vy][vx]) == ' ') {
						continue global
					}
				}
			}
			listOfPotentiallyAppropriateCoords = append(listOfPotentiallyAppropriateCoords, []int{x, y})
		}
	}
	if len(listOfPotentiallyAppropriateCoords) == 0 {
		return nil
	}
	return &listOfPotentiallyAppropriateCoords
}

func (r *RBR) pickJunctionTileForVault(rx, ry int, vaultstrs *[]string, deadendOnly bool) (int, int) {
	listOfAppropriateCoords := make([][]int, 0)
	h, w := len(*vaultstrs), len((*vaultstrs)[0])
	for x := rx; x < rx+w; x++ {
		if rune((*vaultstrs)[0][x-rx]) == '.' && r.isTileSuitableForJunction(x, ry-1, deadendOnly) {
			listOfAppropriateCoords = append(listOfAppropriateCoords, []int{x, ry - 1})
		}
		if rune((*vaultstrs)[h-1][x-rx]) == '.' && r.isTileSuitableForJunction(x, ry+h, deadendOnly) {
			listOfAppropriateCoords = append(listOfAppropriateCoords, []int{x, ry + h})
		}
	}
	for y := ry; y < ry+h; y++ {
		if rune((*vaultstrs)[y-ry][0]) == '.' && r.isTileSuitableForJunction(rx-1, y, deadendOnly) {
			listOfAppropriateCoords = append(listOfAppropriateCoords, []int{rx - 1, y})
		}
		if rune((*vaultstrs)[y-ry][w-1]) == '.' && r.isTileSuitableForJunction(rx+w, y, deadendOnly) {
			listOfAppropriateCoords = append(listOfAppropriateCoords, []int{rx + w, y})
		}
	}

	if len(listOfAppropriateCoords) == 0 {
		return -1, -1
	}
	coord := listOfAppropriateCoords[rnd.Rand(len(listOfAppropriateCoords))]
	return coord[0], coord[1]
}

func (r *RBR) placeRoomvaultByPicking(roomId int, secArea int16, deadendOnly bool) bool {
	placeFound := false
	tries := 0
	maxtries := 1 // Not needed? // r.MAX_RSIZE * r.MAX_RSIZE - r.MIN_RSIZE*r.MIN_RSIZE
finding_place:
	for tries < maxtries {
		tries++
		vaultstrs := r.getRandomRoomvault().getStrings()
		roomH, roomW := len(*vaultstrs), len((*vaultstrs)[0])
		coordsList := r.pickListOfCoordinatesForRoomToBeFit(roomW, roomH)
		// coordsList := r.pickListOfCoordinatesForGivenVaultToBeFit(vaultstrs)
		if coordsList == nil {
			continue finding_place
		}
		selectedCoordIndex := rnd.Rand(len(*coordsList))
		currCoordIndex := selectedCoordIndex
	trying_coords:
		for {
			x, y := (*coordsList)[currCoordIndex][0], (*coordsList)[currCoordIndex][1]
			jx, jy := r.pickJunctionTileForVault(x, y, vaultstrs, deadendOnly)
			if jx != -1 && jy != -1 {
				r.tiles[jx][jy].TileType = TDOOR
				placeFound = true
				r.tryPlaceVaultAtCoords(vaultstrs, x, y, roomId, secArea)
				// r.setRoomIdForTilesRectangle(x, y, roomW, roomH, roomId)
				break finding_place
			}
			currCoordIndex = (currCoordIndex + 1) % len(*coordsList)
			if currCoordIndex == selectedCoordIndex {
				break trying_coords
			}
		}
	}
	return placeFound
}
