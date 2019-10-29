package RBR_generator

import "github.com/sidav/golibrl/random/additive_random"

type RBR struct {
	tiles                    [][]tile
	mapw, maph               int
	MIN_CLENGTH, MAX_CLENGTH int
	MIN_RSIZE, MAX_RSIZE     int
	MINROOMS, MINCORRS       int
	PLACEMENT_TRIES_LIMIT int 
}

var rnd additive_random.FibRandom

func (r *RBR) Init(w, h int) {
	rnd = additive_random.FibRandom{}
	rnd.InitBySeed(-1)
	r.tiles = make([][]tile, w)
	for row := range r.tiles {
		r.tiles[row] = make([]tile, h)
	}
	r.mapw = w
	r.maph = h

	r.MIN_CLENGTH = 2
	r.MAX_CLENGTH = r.mapw - 2 
	r.MIN_RSIZE = 3
	r.MAX_RSIZE = 10 // r.mapw / 10


	// r.MINROOMS = 30
	// r.MINCORRS = 50

	mapArea := r.mapw * r.maph
	maxRoomArea := r.MAX_RSIZE*r.MAX_RSIZE
	minRoomArea := r.MIN_RSIZE*r.MIN_RSIZE
	meanRoomArea := (3*maxRoomArea+minRoomArea)/4
	r.MINROOMS = mapArea / (3* meanRoomArea / 2)
	mapArea -= r.MINROOMS * meanRoomArea
	r.MINCORRS = mapArea / (r.MIN_CLENGTH*10)

	r.PLACEMENT_TRIES_LIMIT = (r.MINROOMS + r.MINCORRS) * 100
}

func (r *RBR) Generate() {

	roomsPlaced, corrsPlaced := r.placeInitialLayout()

	currLoop := 0
	digged := false 
	for (roomsPlaced < r.MINROOMS || corrsPlaced < r.MINCORRS) && currLoop < r.PLACEMENT_TRIES_LIMIT {
		placementFromX, placementfromY := 0, 0 
		placementToX, placementToY := r.mapw, r.maph
		
		roomsRemaining := r.MINROOMS - roomsPlaced
		corrsRemaining := r.MINCORRS - corrsPlaced
		placeRoom := rnd.RandInRange(1, roomsRemaining+corrsRemaining) > corrsRemaining
		if !placeRoom {
			placementFromX += r.MAX_RSIZE/2
			placementfromY += r.MAX_RSIZE/2
			placementToX -= r.MAX_RSIZE/2
			placementToY -= r.MAX_RSIZE/2
		}

		placeOnDeadendOnly := rnd.RandInRange(0, 2) != 0  
		x, y := r.pickJunctionTile(placementFromX, placementfromY, placementToX, placementToY, placeOnDeadendOnly)
		if x == -1 && y == -1 {
			placeOnDeadendOnly = false 
			x, y = r.pickJunctionTile(placementFromX, placementfromY, placementToX, placementToY, false)
		}

		if placeRoom {
			// digged = r.placeRoomFromJunction(x, y, roomsPlaced+1)
			digged = r.placeRoomByPicking(roomsPlaced+1, false)
			if digged {
				roomsPlaced++
			}
		} else {
			forceNotDeadendCorridor := corrsPlaced > r.MINCORRS/4 || roomsPlaced > r.MINROOMS/2
			digged = r.placeCorridorFrom(x, y, forceNotDeadendCorridor)
			if digged {
				corrsPlaced++
			}
		}
		currLoop++
	}
	r.placeRandomDoors(rnd.Rand(r.MINROOMS/5))
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
