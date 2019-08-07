package math

// INTEGER VECTORS
func RotateIntCoords90Degrees(x, y int, counterClockwise bool) (int, int) {
	dir := 1
	if counterClockwise {
		dir = -1
	}
	return -dir*y, dir*x
}

func RotateIntCoords45Degrees(x, y int, counterClockwise bool) (int, int) {
	dir := 1
	if counterClockwise {
		dir = -1
	}
	// no sine/cosine... The code may look awful.
	newvx := x - (dir* y)
	newvy := dir*x + y
	if newvx > 0 {
		newvx = 1
	}
	if newvx < 0 {
		newvx = -1
	}
	if newvy > 0 {
		newvy = 1
	}
	if newvy < 0 {
		newvy = -1
	}
	return newvx, newvy
}
