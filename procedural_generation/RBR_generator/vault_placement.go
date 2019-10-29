package RBR_generator

func (r *RBR) tryPlaceVaultAtCoords(vault *[]string, x, y int) {
	for column := 0; column < len(*vault); column++ {
		for row := 0; row < len((*vault)[column]); row++ {
			symbol := rune((*vault)[column][row])
			ttype := vaultSymbolToTileType(symbol)
			r.tiles[x+row][y+column].tiletype = ttype
		}
	}
}

func (r *RBR) pickListOfCoordinatesForVaultToBeFit(w, h int) *[][]int {
	listOfPotentiallyAppropriateCoords := make([][]int, 0)
	for x := 2; x+w < r.mapw-1; x++ {
		for y := 2; y+h < r.maph-1; y++ {
			if r.isSpaceOfGivenType(x, y, w, h, 1, TFLOOR) {
				listOfPotentiallyAppropriateCoords = append(listOfPotentiallyAppropriateCoords, []int{x, y})
			}
		}
	}
	if len(listOfPotentiallyAppropriateCoords) == 0 {
		return nil
	}
	return &listOfPotentiallyAppropriateCoords
}

func (r *RBR) placeRandomVault() {
	tries := 0
	for tries < len(vaults) {
		vaultStrs := vaults[rnd.Rand(len(vaults))].getStrings()
		w, h := len(*vaultStrs), len((*vaultStrs)[0])
		coordsList := r.pickListOfCoordinatesForVaultToBeFit(w, h)
		if coordsList == nil {
			tries++
			continue
		}
		r.tiles[0][0].tiletype = TDOOR
		coords := (*coordsList)[rnd.Rand(len(*coordsList))]
		r.tryPlaceVaultAtCoords(vaultStrs, coords[0], coords[1])
		break
	}
}
