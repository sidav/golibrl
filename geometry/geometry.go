package geometry

func AreCoordsInRect(x, y, rx, ry, w, h int) bool {
	return x >= rx && x < rx+w && y >= ry && y < ry+h
}

func AreCoordsInRange(fx, fy, tx, ty, r int) bool { // border including.
	// uses more wide circle (like in Bresenham's circle) than the real geometric one.
	// It is much more handy for spaces with discrete coords (cells).
	realSqDistanceAndSqRadiusDiff := (fx-tx)*(fx-tx) + (fy-ty)*(fy-ty) - r*r
	return realSqDistanceAndSqRadiusDiff < r
}

func AreCoordsInRangeFromRect(fx, fy, tx, ty, w, h, r int) bool { // considering ANY of the tiles in the rect.
	return AreRectsInRange(fx, fy, 1, 1, tx, ty, w, h, r)
}

func AreRectsInRange(x1, y1, w1, h1, x2, y2, w2, h2, r int) bool {
	// all -1's are beacuse of TILED geometry
	x1b := x1+w1-1
	x2b := x2+w2-1
	y1b := y1+h1-1
	y2b := y2+h2-1
	
	left := x2b < x1
	right := x1b < x2
	bottom := y1b < y2
	top := y2b < y1
	if top && left {
		return AreCoordsInRange(x1, y1, x2b, y2b, r) // dist((x1, y1b), (x2b, y2))
	}
	if left && bottom {
		return AreCoordsInRange(x1, y1b, x2b, y2, r)
	}
	if bottom && right {
		return AreCoordsInRange(x1b, y1b, x2, y2, r)
	}
	if right && top {
		return AreCoordsInRange(x1b, y1, x2, y2b, r)
	}
	if left {
		return x1 - x2b <= r
	}
	if right {
		return x2 - x1b <= r
	}
	if bottom {
		return y2 - y1b <= r
	}
	if top {
		return y1 - y2b <= r
	}
	return true // intersect detected
}

func GetCellNearestToRectFrom(rx, ry, w, h, fx, fy int) (int, int) { // returns a cell closest to rectangle
	left := fx < rx
	right := fx > rx+w-1
	bottom := fy > ry+h-1
	top := fy < ry
	if top && left {
		return rx-1, ry-1
	}
	if left && bottom {
		return rx-1, ry+h
	}
	if bottom && right {
		return rx+w, ry+h
	}
	if right && top {
		return rx+w, ry-1
	}
	if left {
		return rx-1, fy
	}
	if right {
		return rx+w, fy
	}
	if bottom {
		return fx, ry+h
	}
	if top {
		return fx, ry-1
	}
	return fx, fy // intersect
}

func getSqDistanceBetween(x1, y1, x2, y2 int) int {
	return (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
}

func AreTwoCellRectsOverlapping(x1, y1, w1, h1, x2, y2, w2, h2 int) bool {
	// WARNING:
	// ALL "-1"s HERE ARE BECAUSE OF WE ARE IN CELLS SPACE
	// I.E. A SINGLE CELL IS 1x1 RECTANGLE
	// SO RECTS (0, 0, 1x1) AND (1, 0, 1x1) ARE NOT OVERLAPPING IN THIS SPACE (BUT SHOULD IN EUCLIDEAN OF COURSE)
	right1 := x1 + w1 - 1
	bot1 := y1 + h1 - 1
	right2 := x2 + w2 - 1
	bot2 := y2 + h2 - 1
	return !(x2 > right1 ||
		right2 < x1 ||
		y2 > bot1 ||
		bot2 < y1)
}

//func AreCircleAndRectOverlapping(cx, cy, r, rx, ry, w, h int) bool {
//	// topleft: rx, ry
//	// topright: rx+w, ry
//	// downleft: rx, ry+h
//	// downright: rx+w, ry+h
//	// will work bad (bad case example: it won't detect veeeeeeeery wide rectangle and small circle intersection if no corners are in the circle), but suitable for the game, I hope...
//	return AreCoordsInRange(rx, ry, cx, cy, r) || AreCoordsInRange(rx+w, ry, cx, cy, r) || AreCoordsInRange(rx, ry+h, cx, cy, r) || AreCoordsInRange(rx+w, ry+h, cx, cy, r)
//}
