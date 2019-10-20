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

	r.MIN_CLENGTH = 3
	r.MAX_CLENGTH = r.mapw / 2
	r.MIN_RSIZE = 3
	r.MAX_RSIZE = r.mapw / 8


	r.MINROOMS = 30
	r.MINCORRS = 40

	// mapArea := r.mapw * r.maph
	// maxRoomArea := r.MAX_RSIZE*r.MAX_RSIZE
	// r.MINROOMS = mapArea / (3 * maxRoomArea / 2)
	// mapArea -= r.MINROOMS * maxRoomArea
	// r.MINCORRS = mapArea / (3*r.MAX_CLENGTH)

	r.PLACEMENT_TRIES_LIMIT = (r.MINROOMS + r.MINCORRS) * 10
}

func (r *RBR) Generate() {
	// place initial room
	digged := false
	for !digged {
		x := rnd.RandInRange(0, r.mapw)
		y := rnd.RandInRange(0, r.maph)
		w := rnd.RandInRange(5, 10)
		h := rnd.RandInRange(5, 10)
		digged = r.digRoomIfPossible(x, y, w, h, 1)
	}

	roomsRemaining := r.MINROOMS - 1
	corrsRemaining := r.MINCORRS
	currLoop := 0

	for (roomsRemaining != 0 || corrsRemaining != 0) && currLoop < r.PLACEMENT_TRIES_LIMIT {
		placeOnDeadendOnly := rnd.RandInRange(0, 1) == 1 
		x, y := r.pickJunctionTile(placeOnDeadendOnly)
		if x == -1 && y == -1 {
			x, y = r.pickJunctionTile(false)
		}
		placeRoom := rnd.RandInRange(1, roomsRemaining+corrsRemaining) > corrsRemaining
		if placeRoom {
			digged = r.placeRoomFromJunction(x, y)
			if digged {
				roomsRemaining--
			}
		} else {
			digged = r.placeCorridorFrom(x, y)
			if digged {
				corrsRemaining--
			}
		}
		currLoop++
	}
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
