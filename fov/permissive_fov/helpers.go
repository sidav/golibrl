package permissive_fov

import "math"

func getDistance(x, y int) int {
	return int(math.Sqrt(float64(x*x + y*y)))
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

