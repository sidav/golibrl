package ai_helpers

func sqdist(x1, y1, x2, y2 int) int {
	return (x2-x1)*(x2-x1)+(y2-y1)*(y2-y1)
}

// Deprecated
// Naive method
func FindCoordsByConditionAndClosestFrom2(condition func(int, int) bool, fromX, fromY, maxDist int) (int, int, bool) {
	if condition(fromX, fromY) {
		return fromX, fromY, true
	}
	var candidateX, candidateY int
	candidateFound := false
	for x := fromX-maxDist; x <= fromX+maxDist; x++ {
		for y := fromY-maxDist; y <= fromY+maxDist; y++ {
			if condition(x, y) {
				if !candidateFound || sqdist(fromX, fromY, x, y) < sqdist(fromX, fromY, candidateX, candidateY) {
					candidateX = x
					candidateY = y
					candidateFound = true
				}
			}
		}
	}
	return candidateX, candidateY, candidateFound
}

func FindCoordsByConditionAndClosestFrom(condition func(int, int) bool, fromX, fromY, maxDist int) (int, int, bool) {
	x, y := fromX, fromY
	if condition(x, y) {
		return x, y, true
	}

	y += 1
	dirX, dirY := 1, 0
	currRadius := 1 // euclidean
	roundStartX, roundStartY := x, y
	var candidateX, candidateY, candidateSqDist int
	candidateFound := false

	for {
		// check condition and set candidate if it is closer than previous
		if condition(x, y) {
			if !candidateFound || sqdist(fromX, fromY, x, y) < candidateSqDist {
				candidateX = x
				candidateY = y
				candidateSqDist = sqdist(fromX, fromY, candidateX, candidateY)
				candidateFound = true
			}
		}
		// check if need to change direction
		if (x + dirX) - fromX > currRadius || fromX - (x+dirX) > currRadius || y + dirY - fromY > currRadius || fromY - (y+dirY) > currRadius {
			dirX, dirY = dirY, -dirX // rotate 90 degs clockwise
		}
		x += dirX
		y += dirY
		if x == roundStartX && y == roundStartY { // increase radius and go on
			currRadius += 1
			// check if candidate was found, and any subsequent candidates won't be any closer
			if candidateFound && candidateSqDist <= currRadius*currRadius {
				return candidateX, candidateY, true
			}
			if currRadius > maxDist {
				return -1, -1, false
			}
			y += 1
			roundStartX, roundStartY = x, y
		}
	}
}

