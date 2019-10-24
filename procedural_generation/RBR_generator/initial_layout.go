package RBR_generator

// returns number of rooms and corridors placed.
func (r *RBR) placeInitialLayout() (int, int) {
	// r.placeInitialSingleRoom()
	corridorRings := rnd.RandInRange(1, 5)
	for i := 0; i < corridorRings; i++ {
		r.placeInitialCorridorRing()
	}
	return 0, 4
}

func (r *RBR) placeInitialSingleRoom() {
	digged := false
	for !digged {

		x := rnd.RandInRange(r.MAX_RSIZE, r.mapw-r.MAX_RSIZE-r.MIN_RSIZE)
		y := rnd.RandInRange(r.MAX_RSIZE, r.maph-r.MAX_RSIZE-r.MIN_RSIZE)

		w := rnd.RandInRange(r.MIN_RSIZE, r.MAX_RSIZE)
		h := rnd.RandInRange(r.MIN_RSIZE, r.MAX_RSIZE)
		r.digSpace(x, y, w, h, 1)
		digged = true
	}
}

func (r *RBR) placeInitialCorridorRing() {
	x, y := rnd.RandInRange(1, r.mapw-r.MIN_RSIZE-1), rnd.RandInRange(1, r.maph-r.MIN_RSIZE-1)
	w, h := rnd.RandInRange(r.MIN_RSIZE, r.mapw-x-2), rnd.RandInRange(r.MIN_RSIZE, r.maph-y-2)

	for cx := x; cx <= x+w; cx++ {
		for cy := y; cy <= y+h; cy++ {
			r.digSpace(cx, y, 1, 1, 0)
			r.digSpace(cx, y+h, 1, 1, 0)
			r.digSpace(x, cy, 1, 1, 0)
			r.digSpace(x+w, cy, 1, 1, 0)
		}
	}
}
