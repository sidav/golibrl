package RBR_generator

// returns number of rooms and corridors placed.
func (r *RBR) placeInitialLayout() (int, int) {
	layout := rnd.RandInRange(0, 3)
	
	rooms, corrs := 0, 0
	switch layout {
	case 0:
		corridorRings := rnd.RandInRange(1, 5)
		rooms, corrs = r.placeInitialCorridorRings(corridorRings)
	case 1:
		rooms, corrs = r.placeInitialTwoInterconnectedRooms()
	case 2:
		rooms, corrs = r.placeInitialFourInterconnectedRooms()
	default:
		rooms, corrs = r.placeInitialLargeRoom()
	}
	return rooms, corrs
}

func (r *RBR) placeInitialLargeRoom() (int, int) {
	digged := false
	for !digged {

		x := rnd.RandInRange(3, r.mapw/4)
		y := rnd.RandInRange(r.maph/4, r.maph/2)
		w := r.mapw - x - x - 2 // WARNING: violates room size constraints!
		h := rnd.RandInRange(r.maph/5, r.maph/2-1)
		r.digSpace(x, y, w, h, 1)
		digged = true
	}
	return 1, 0
}

func (r *RBR) placeInitialTwoInterconnectedRooms() (int, int) {
	w1 := rnd.RandInRange(r.MAX_RSIZE/2, r.MAX_RSIZE)
	w2 := w1 // rnd.RandInRange(r.MIN_RSIZE, r.MAX_RSIZE)
	x1 := rnd.RandInRange(1, r.mapw/3)
	x2 := r.mapw - x1 - w2 // rnd.RandInRange(r.mapw/2, r.mapw-w2)
	h := rnd.RandInRange(r.MAX_RSIZE/2, r.MAX_RSIZE)
	y := rnd.RandInRange(r.maph/4, 3*r.maph/4-h)
	r.digSpace(x1, y, w1, h, 0)
	r.digSpace(x2, y, w2, h, 1)
	corridors := rnd.RandInRange(1, h/3)
	if corridors < 1 {
		corridors = 1
	}
	for i := 0; i < corridors; i++ {
		corrY := rnd.RandInRange(y, y+h-1)
		r.digSpace(x1+w1, corrY, x2-x1-w1, 1, 0)
	}
	return 2, corridors
}

func (r *RBR) placeInitialFourInterconnectedRooms() (int, int) {
	w := rnd.RandInRange(r.MAX_RSIZE/2, r.MAX_RSIZE)
	x1 := rnd.RandInRange(1, r.mapw/3)
	x2 := r.mapw - x1 - w 
	h := rnd.RandInRange(r.MAX_RSIZE/2, r.MAX_RSIZE)
	y1 := rnd.RandInRange(1, r.maph/3)
	y2 := r.maph - y1 - h 
	r.digSpace(x1, y1, w, h, 1)
	r.digSpace(x2, y1, w, h, 2)
	r.digSpace(x1, y2, w, h, 3)
	r.digSpace(x2, y2, w, h, 4)
	corridors := rnd.RandInRange(1, h/3)
	if corridors < 1 {
		corridors = 1
	}
	for i := 0; i < corridors; i++ {
		corrY := rnd.RandInRange(y1, y1+h-1)
		r.digSpace(x1+w, corrY, x2-x1-w, 1, 0)
		r.digSpace(x1+w, r.maph-corrY-1, x2-x1-w, 1, 0)
		corrX := rnd.RandInRange(x1, x1+w-1)
		r.digSpace(corrX, y1+h, 1, y2-y1-h, 0)
		r.digSpace(r.mapw-corrX-1, y1+h, 1, y2-y1-h, 0)
	}
	return 4, 4 * corridors
}

func (r *RBR) placeInitialCorridorRings(number int) (int, int) {
	for i := 0; i < number; i++ {
		x, y := rnd.RandInRange(1, r.mapw-r.MIN_RSIZE-1), rnd.RandInRange(1, r.maph-r.MIN_RSIZE-1)
		w, h := rnd.RandInRange(r.MIN_RSIZE, r.mapw-x-2), rnd.RandInRange(r.MIN_RSIZE, r.maph-y-2)
		if i == 0 {
			// The first corridor ring SHOULD cover the most of the map to maximize map space usage.
			x, y = rnd.RandInRange(1, r.MAX_RSIZE), rnd.RandInRange(1, r.MAX_RSIZE)
			w, h = rnd.RandInRange(r.mapw/2, r.mapw-x-2), rnd.RandInRange(r.maph/2, r.maph-y-2)
		}

		for cx := x; cx <= x+w; cx++ {
			for cy := y; cy <= y+h; cy++ {
				r.digSpace(cx, y, 1, 1, 0)
				r.digSpace(cx, y+h, 1, 1, 0)
				r.digSpace(x, cy, 1, 1, 0)
				r.digSpace(x+w, cy, 1, 1, 0)
			}
		}
	}
	return 0, number * 4
}
