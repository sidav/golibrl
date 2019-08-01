package geometry

import "math"

func AreCoordsInSector(fromX, fromY, lookX, lookY, targetX, targetY, sectorAngle int) bool {
	if fromX == targetX && fromY == targetY {
		return true
	}
	if sectorAngle > 180 { // doesn't work
		return !AreCoordsInSector(fromX, fromY, lookX, lookY, targetX, targetY, 360 - sectorAngle)
	}
	halfSectorAngle := math.Pi * (float64(sectorAngle) / 2) / 180  // in radians
	centX, centY := targetX-fromX, targetY-fromY
	sectorCenterAngle := math.Atan2(float64(lookY), float64(lookX))
	angleToCoords := math.Atan2(float64(centY), float64(centX))
	if centX < 0 && centY < 0 && lookY >= 0 {
		angleToCoords += 2 * math.Pi
	}
	return sectorCenterAngle - halfSectorAngle <= angleToCoords && sectorCenterAngle + halfSectorAngle >= angleToCoords
}
