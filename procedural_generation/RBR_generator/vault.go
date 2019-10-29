package RBR_generator

import "os"
import "bufio"
import "strings"

var vaults [][]string

func readVaultsFromFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	vaultLines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.Contains(line, "//") {
			if len(vaultLines) > 0 {
				vaults = append(vaults, vaultLines)
				vaultLines = make([]string, 0)
			}
		} else {
			vaultLines = append(vaultLines, line)
		}
	}
	if len(vaultLines) > 0 {
		vaults = append(vaults, vaultLines)
	}
}

func vaultSymbolToTileType(symbol rune) byte {
	switch symbol {
	case '#':
		return TWALL
	case '+':
		return TDOOR
	case '.':
		return TFLOOR
	default:
		return TDOOR
	}
}

func (r *RBR) tryPlaceVaultAtCoords(vault *[]string, x, y int) {
	for row := 0; row < len(*vault); row++ {
		for col := 0; col < len((*vault)[row]); col++ {
			symbol := rune((*vault)[row][col])
			ttype := vaultSymbolToTileType(symbol)
			r.tiles[x+row][y+col].tiletype = ttype
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
		vault := &vaults[rnd.Rand(len(vaults))]
		h, w := len(*vault), len((*vault)[0])
		coordsList := r.pickListOfCoordinatesForVaultToBeFit(w+1, h+1)
		if coordsList == nil {
			tries++
			continue
		}
		r.tiles[0][0].tiletype = TDOOR
		coords := (*coordsList)[rnd.Rand(len(*coordsList))]
		r.tryPlaceVaultAtCoords(vault, coords[0]+1, coords[1]+1)
		break
	}
}
