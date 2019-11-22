package RBR_generator

import "github.com/sidav/golibrl/random/additive_random"

type RBR struct {
	tiles                                               [][]tile
	mapw, maph                                          int
	MIN_CLENGTH, MAX_CLENGTH                            int
	MIN_RSIZE, MAX_RSIZE, ROOM_SIZE_BIAS                int
	MINROOMS, MINCORRS                                  int
	PLACEMENT_TRIES_LIMIT                               int
	VAULTS_NUM                                          int
	numPlacedRooms, numPlacedVaults, numPlacedCorridors int
	vaults, roomvaults                                  []*vault
	NUM_SEC_AREAS                                       int
}

var rnd additive_random.FibRandom

func (r *RBR) Init(w, h, secareas int, vaultsFilePath, roomvaultsFilePath string) {
	rnd = additive_random.FibRandom{}

	if vaultsFilePath != "" {
		r.readVaultsFromFile("procedural_generation/RBR_generator/vaults.txt") // TODO: custom vaults file.
	}
	if roomvaultsFilePath != "" {
		r.readRoomVaultsFromFile("procedural_generation/RBR_generator/roomvaults.txt")
	}

	rnd.InitBySeed(-1)
	r.tiles = make([][]tile, w)
	for row := range r.tiles {
		r.tiles[row] = make([]tile, h)
	}
	r.mapw = w
	r.maph = h
	r.NUM_SEC_AREAS = secareas

	// TODO: make these configurable.
	r.MIN_CLENGTH = 2
	r.MAX_CLENGTH = r.mapw - 2
	r.MIN_RSIZE = 3
	r.MAX_RSIZE = (r.mapw - 2) / 7
	r.ROOM_SIZE_BIAS = r.MAX_RSIZE / 2
	r.VAULTS_NUM = len(r.vaults)

	// r.MINROOMS = 30
	// r.MINCORRS = 50

	mapArea := r.mapw * r.maph
	// maxRoomArea := r.MAX_RSIZE * r.MAX_RSIZE
	// minRoomArea := r.MIN_RSIZE * r.MIN_RSIZE
	meanRoomArea := (r.ROOM_SIZE_BIAS * r.ROOM_SIZE_BIAS)
	r.MINROOMS = mapArea / (3 * meanRoomArea / 2)
	mapArea -= r.MINROOMS * meanRoomArea
	r.MINCORRS = mapArea / (r.MIN_CLENGTH * 20)
	// r.MINCORRS = 0

	r.PLACEMENT_TRIES_LIMIT = (r.MINROOMS + r.MINCORRS) * 10
}

func (r *RBR) Generate() {

	r.numPlacedRooms, r.numPlacedCorridors = r.placeInitialLayout()

	currLoop := 0
	digged := false
	increaseSecAreaEach := r.MINROOMS / r.NUM_SEC_AREAS
	increaseSecAreaEach = rnd.RandInRange(increaseSecAreaEach-increaseSecAreaEach/2, increaseSecAreaEach+increaseSecAreaEach/2)
	var currSecArea int16 = 0

	for (r.numPlacedRooms < r.MINROOMS || r.numPlacedCorridors < r.MINCORRS) && currLoop < r.PLACEMENT_TRIES_LIMIT {
		if r.numPlacedRooms%increaseSecAreaEach == 0 && int(currSecArea) < r.NUM_SEC_AREAS-1 {
			currSecArea++
		}
		placementFromX, placementfromY := 0, 0
		placementToX, placementToY := r.mapw, r.maph

		roomsRemaining := r.MINROOMS - r.numPlacedRooms
		corrsRemaining := r.MINCORRS - r.numPlacedCorridors
		placeRoom := rnd.RandInRange(1, roomsRemaining+corrsRemaining) > corrsRemaining
		if !placeRoom {
			// change placement bounds ('cause it's corridor and we wanna reduce deadends) 
			placementFromX += r.ROOM_SIZE_BIAS
			placementfromY += r.ROOM_SIZE_BIAS
			placementToX -= r.ROOM_SIZE_BIAS
			placementToY -= r.ROOM_SIZE_BIAS
		}

		placeOnDeadendOnly := rnd.RandInRange(0, 2) != 0
		x, y := r.pickJunctionTile(placementFromX, placementfromY, placementToX, placementToY, placeOnDeadendOnly)
		if x == -1 && y == -1 {
			placeOnDeadendOnly = false
			x, y = r.pickJunctionTile(placementFromX, placementfromY, placementToX, placementToY, false)
		}

		if placeRoom {
			digged = false
			if rnd.OneChanceFrom(2) || r.numPlacedVaults >= r.VAULTS_NUM {
				vaultNeeded := rnd.OneChanceFrom(3)
				digged = r.placeRoomByPicking(r.numPlacedRooms+1, currSecArea, false, vaultNeeded)
			} else {
				digged = r.placeRoomvaultByPicking(r.numPlacedRooms+1, currSecArea, false)
			}
			if digged {
				r.numPlacedRooms++
			}
		} else {
			forceNotDeadendCorridor := r.numPlacedCorridors > r.MINCORRS/4 || r.numPlacedRooms > r.MINROOMS/2
			digged = r.placeCorridorFrom(x, y, forceNotDeadendCorridor)
			if digged {
				r.numPlacedCorridors++
			}
		}
		currLoop++
	}
	r.placeRandomDoors(rnd.Rand(r.MINROOMS / 5))
	r.finalizeDoorsSecArea()
	r.placeStairs(1, 2, true)
	// for i := r.numPlacedVaults; i < r.VAULTS_NUM; i++ {
	// 	r.placeRandomVault()
	// }
}

///////////////////////////////////////////////////////////////////
func (rbr *RBR) GetMapChars() *[][]rune {
	runearr := make([][]rune, rbr.mapw)
	for row := range runearr {
		runearr[row] = make([]rune, rbr.maph)
	}
	for x := 0; x < rbr.mapw; x++ {
		for y := 0; y < rbr.maph; y++ {
			runearr[x][y] = rbr.tiles[x][y].toRune()
		}
	}
	return &runearr
}
