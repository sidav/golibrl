package RBR_generator

func (r *RBR) tryPlaceVaultAtCoords(vaultString *[]string, x, y int) {
	for column := 0; column < len(*vaultString); column++ {
		for row := 0; row < len((*vaultString)[column]); row++ {
			symbol := rune((*vaultString)[column][row])
			ttype := vaultSymbolToTileType(symbol)
			r.tiles[x+row][y+column].tiletype = ttype
			r.numPlacedVaults++
		}
	}
}

func (r *RBR) tryPlaceVaultOfGivenSizeAtCoords(x, y, w, h int) {
	vaultsOfSize := make([]*vault, 0)
	for _, v := range vaults {
		if v.isOfSize(w, h) {
			vaultsOfSize = append(vaultsOfSize, v)
		}
	}
	if len(vaultsOfSize) == 0 {
		return
	}
	r.tiles[0][0].tiletype = TDOOR
	vlt := vaultsOfSize[rnd.Rand(len(vaultsOfSize))]
	r.tryPlaceVaultAtCoords(vlt.getStringsIfFitInSize(w, h), x, y)
}

func (r *RBR) pickListOfCoordinatesForVaultToBeFit(w, h int) *[][]int {
	listOfPotentiallyAppropriateCoords := make([][]int, 0)
	for x := 2; x+w < r.mapw-1; x++ {
		for y := 2; y+h < r.maph-1; y++ {
			if r.isSpaceOfGivenType(x, y, w, h, 1, TFLOOR) && !r.isSpaceOfGivenType(x, y, w, h, 2, TFLOOR) {
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
		h, w := len(*vaultStrs), len((*vaultStrs)[0])
		coordsList := r.pickListOfCoordinatesForVaultToBeFit(w, h)
		if coordsList == nil {
			tries++
			continue
		}
		coords := (*coordsList)[rnd.Rand(len(*coordsList))]
		r.tryPlaceVaultAtCoords(vaultStrs, coords[0], coords[1])
		break
	}
}
