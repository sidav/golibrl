package RBR_generator

import "github.com/sidav/golibrl/random/additive_random"

type RBR struct {
	tiles                    [][]tile
	mapw, maph               int
	MIN_CLENGTH, MAX_CLENGTH int
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
	r.MAX_CLENGTH = 10
}

func (r *RBR) Generate() {

	for room := 0; room < 1; room++ {
		success := false
		for !success {
			x := rnd.RandInRange(0, r.mapw)
			y := rnd.RandInRange(0, r.maph)
			w := rnd.RandInRange(5, 10)
			h := rnd.RandInRange(5, 10)
			success = r.digRoomIfPossible(x, y, w, h, 1)
		}
	}
	for crrdr := 0; crrdr < 100; crrdr++ {
		x, y := r.pickJunctionTile()
		digged := r.placeCorridorFrom(x, y)
		if !digged {
			crrdr--
			continue
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
