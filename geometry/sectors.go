package geometry

import (
	"math"
)

func AreCoordsInSector(x, y, sectorOriginX, sectorOriginY, sectorDirX, sectorDirY, sectorAngle int) bool {
	inverse := false
	if sectorOriginX == x && sectorOriginY == y {
		return true
	}
	if sectorAngle > 180 {
		if sectorAngle >= 360 {
			return true
		}
		sectorAngle = 360 - sectorAngle
		inverse = true
		sectorDirX = -sectorDirX
		sectorDirY = -sectorDirY
	}
	halfSectorAngle := math.Pi * (float64(sectorAngle) / 2) / 180 // in radians
	centX, centY := x-sectorOriginX, y-sectorOriginY
	sectorCenterAngle := math.Atan2(float64(sectorDirY), float64(sectorDirX))
	angleToCoords := math.Atan2(float64(centY), float64(centX))
	if centX < 0 && centY < 0 && sectorDirY >= 0 {
		angleToCoords += 2 * math.Pi
	}
	liesWithin := sectorCenterAngle-halfSectorAngle <= angleToCoords && sectorCenterAngle+halfSectorAngle >= angleToCoords
	if inverse {
		return !liesWithin
	}
	return liesWithin
}
