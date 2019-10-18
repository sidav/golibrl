package RBR_generator

import "github.com/sidav/golibrl/random/additive_random"

type RBR struct {
	tiles                    [][]tile
	mapw, maph               int
	MIN_CLENGTH, MAX_CLENGTH int
	MIN_RSIZE, MAX_RSIZE     int
	MINROOMS, MINCORRS       int
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
	r.MAX_CLENGTH = 10//r.mapw
	r.MIN_RSIZE = 3
	r.MAX_RSIZE = 10


	r.MINROOMS = 20
	r.MINCORRS = 1

	mapArea := r.mapw * r.maph
	maxRoomArea := r.MAX_RSIZE*r.MAX_RSIZE
	r.MINROOMS = mapArea / (3 * maxRoomArea / 2)
	mapArea -= r.MINROOMS * maxRoomArea
	r.MINCORRS = mapArea / (3*r.MAX_CLENGTH)
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

	for roomsRemaining != 0 || corrsRemaining != 0 {
		x, y := r.pickJunctionTile()
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
