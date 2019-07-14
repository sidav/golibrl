package mill_fov

import "github.com/sidav/golibrl/geometry"

// This algorithm is taken from http://www.adammil.net/blog/v125_roguelike_vision_algorithms.html
// Rewritten in Go by sidav

var (
	mapw, maph, rangeLimit int
	opaque                 *[][]bool
	visible                *[][]bool
)

func SetOpacityMap(o *[][]bool) {
	opaque = o
	mapw = len(*opaque)
	maph = len((*opaque)[0])
}

func emptyVisibilityMap(w, h int) {
	vis := make([][]bool, w)
	for i := range vis {
		vis[i] = make([]bool, h)
	}
	visible = &vis
}

func Fov(fx, fy, rangeLimit int) *[][]bool {
	emptyVisibilityMap(mapw, maph)
	(*visible)[fx][fy] = true
	var octant byte
	for octant =0; octant <8; octant++{
		computeOctant(octant, fx, fy, rangeLimit, 1, &slope{1, 1}, &slope{0, 1})
	}
	return visible
}

func computeOctant(octant byte, fx, fy, rangeLimit, x int, top, bottom *slope) {
	// throughout this function there are references to various parts of tiles. a tile's coordinates refer to its
	// center, and the following diagram shows the parts of the tile and the vectors from the origin that pass through
	// those parts. given a part of a tile with vector u, a vector v passes above it if v > u and below it if v < u
	//    g         center:        Y / X
	// a------b   a top left:      (Y*2+1) / (X*2-1)   i inner top left:      (Y*4+1) / (X*4-1)
	// |  /\  |   b top right:     (Y*2+1) / (X*2+1)   j inner top right:     (Y*4+1) / (X*4+1)
	// |i/__\j|   c bottom left:   (Y*2-1) / (X*2-1)   k inner bottom left:   (Y*4-1) / (X*4-1)
	//e|/|  |\|f  d bottom right:  (Y*2-1) / (X*2+1)   m inner bottom right:  (Y*4-1) / (X*4+1)
	// |\|__|/|   e middle left:   (Y*2) / (X*2-1)
	// |k\  /m|   f middle right:  (Y*2) / (X*2+1)     a-d are the corners of the tile
	// |  \/  |   g top center:    (Y*2+1) / (X*2)     e-h are the corners of the inner (wall) diamond
	// c------d   h bottom center: (Y*2-1) / (X*2)     i-m are the corners of the inner square (1/2 tile width)
	//    h
	for ; x <= rangeLimit; x++ {
		// compute the Y coordinates of the top and bottom of the sector. we maintain that top > bottom
		var topY int
		if top.X == 1 {
			topY = x
		} else {
			topY = ((x*2-1) * top.Y + top.X) / (top.X*2)

			if blocksLight(x, topY, octant, fx, fy) {
				if top.GreaterOrEqual(topY*2+1, x*2) && !blocksLight(x, topY+1, octant, fx, fy) {
					topY++
				}
			} else {
				ax := x*2
				if blocksLight(x+1, topY+1, octant, fx, fy) {
					ax++
				}
				if top.Greater(topY*2+1, ax) {
					topY++
				}
			}
		}
		var bottomY int
		if bottom.Y == 0 {
			bottomY = 0
		} else {
			bottomY = ((x*2-1) * bottom.Y + bottom.X) / (bottom.X*2)
			if bottom.GreaterOrEqual(bottomY*2+1, x*2) && blocksLight(x, bottomY, octant, fx, fy) &&
				!blocksLight(x, bottomY+1, octant, fx, fy) {
					bottomY++
			}
		}
		wasOpaque := -1
		for y:=topY; y >= bottomY; y-- {
			if geometry.AreCoordsInRange(x, y, 0, 0, rangeLimit) { // TODO: range check
			isOpaque := blocksLight(x, y, octant, fx, fy)
			isVisible := isOpaque || ((y != topY || top.Greater(y*4-1, x*4+1)) && (y != bottomY || bottom.Less(y*4+1, x*4-1)))
			if isVisible {
				setVisible(x, y, octant, fx, fy)
			}

			if x != rangeLimit {
				if isOpaque {
					if wasOpaque == 0 {
						nx := x*2
						ny := y*2+1
						if blocksLight(x, y+1, octant, fx, fy) {
							nx--
						}
						if top.Greater(ny, nx) {
							if y == bottomY {
								bottom = &slope{ny, nx}
								break
							} else {
								computeOctant(octant, fx, fy, rangeLimit, x+1, top, &slope{ny, nx})
							}
						} else {
							if y == bottomY {
								return
							}
						}
					}
					wasOpaque = 1
				} else {
					if wasOpaque > 0 {
						nx := x*2
						ny := y*2 + 1
						if blocksLight(x+1, y+1, octant, fx, fy) {
							nx++
						}
						if bottom.GreaterOrEqual(ny, nx) {
							return
						}
						top = &slope{ny, nx}
					}
					wasOpaque = 0
				}
			}

			// TO HERE
			}
		}
		if wasOpaque != 0 {
			break
		}
	}
}

func blocksLight(x, y int, octant byte, fx, fy int) bool {
	nx := fx
	ny := fy
	switch octant {
	case 0: nx += x; ny -= y; break
	case 1: nx += y; ny -= x; break
	case 2: nx -= y; ny -= x; break
	case 3: nx -= x; ny -= y; break
	case 4: nx -= x; ny += y; break
	case 5: nx -= y; ny += x; break
	case 6: nx += y; ny += x; break
	case 7: nx += x; ny += y; break
	}
	if geometry.AreCoordsInRect(nx, ny, 0, 0, mapw, maph) {
		return (*opaque)[nx][ny]
	} else {
		return true
	}
}

func setVisible(x, y int, octant byte, fx, fy int) {
	nx := fx
	ny := fy
	switch octant {
	case 0: nx += x; ny -= y; break
	case 1: nx += y; ny -= x; break
	case 2: nx -= y; ny -= x; break
	case 3: nx -= x; ny -= y; break
	case 4: nx -= x; ny += y; break
	case 5: nx -= y; ny += x; break
	case 6: nx += y; ny += x; break
	case 7: nx += x; ny += y; break
	}
	if geometry.AreCoordsInRect(nx, ny, 0, 0, mapw, maph) {
		(*visible)[nx][ny] = true
	}
}
