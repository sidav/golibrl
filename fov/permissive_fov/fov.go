package permissive_fov

import (
	"fmt"
	"github.com/sidav/golibrl/console"
	"github.com/sidav/golibrl/geometry"
)

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

///

var (
	quadrant, source offset
)

func GetFovMapFrom(x, y, radius int) *[][]bool {
	source = offset{x, y}
	rangeLimit = radius
	emptyVisibilityMap(mapw, maph)
	for q := 0; q < 4; q++ {
		switch q {
		case 0:
			quadrant.x = 1
			quadrant.y = 1
		case 1:
			quadrant.x = -1
			quadrant.y = 1
		case 2:
			quadrant.x = -1
			quadrant.y = -1
		case 3:
			quadrant.x = 1
			quadrant.y = -1
		}
		computeQuadrant()
	}
	return visible
}

func computeQuadrant() {
	var infinity = mapw+maph// int(^uint32(0) >> 1)
	steepBumps := make([]*bump, 0)
	shallowBumps := make([]*bump, 0)
	activeFields := fieldList{}
	activeFields.addToEnd(&field{
		steep:   line{offset{1, 0}, offset{0, maph}},
		shallow: line{offset{0, 1}, offset{mapw, 0}},
	})
	dest := offset{}
	if quadrant.x == 1 && quadrant.y == 1 {
		actIsBlocked(dest)
	}
	for i := 1; i < infinity && activeFields.size > 0; i++ {
		startJ := max(0, i - mapw)
		maxJ := min(i, maph)
		current := activeFields.first
		for j := startJ; j <= maxJ; j++ {
			dest.x = i - j
			dest.y = j
			current = visitSquare(dest, current, &steepBumps, &shallowBumps, &activeFields)
		}
	}
}

func actIsBlocked(pos offset) bool {
	if rangeLimit >= 0 && !geometry.AreCoordsInRange(pos.x, pos.y, 0, 0, rangeLimit) {
		return true
	}
	x := pos.x*quadrant.x + source.x
	y := pos.y*quadrant.y + source.y
	if geometry.AreCoordsInRect(x, y, 0, 0, mapw, maph) {
		(*visible)[x][y] = true
		return (*opaque)[x][y]
	}
	return true // squares out of bounds of the opacity map are considered opaque.
}

func visitSquare(dest offset, currentField *field, steepBumps *[]*bump, shallowBumps *[]*bump, activeFields *fieldList) *field {
	topLeft := offset{dest.x, dest.y + 1}
	bottomRight := offset{dest.x + 1, dest.y}

	for currentField != nil && currentField.steep.isBelowOrContains(bottomRight) {
		// case ABOVE
		// The square is in case 'above'. This means that it is ignored
		// for the currentField. But the steeper fields might need it.
		currentField = currentField.next
	}

	if currentField == nil {
		// The square was in case 'above' for all fields. This means that
		// we no longer care about it or any squares in its diagonal rank.
		return currentField
	}

	if currentField.shallow.isAboveOrContains(topLeft) {
		// case BELOW
		// The shallow line is above the extremity of the square, so that
		// square is ignored.
		return currentField
	}

	// The square is between the lines in some way. This means that we
	// need to visit it and determine whether it is blocked.

	if !actIsBlocked(dest) {
		// We don't care what case might be left, because this square does
		// not obstruct.
		return currentField
	}

	if currentField.shallow.isAbove(bottomRight) &&
		currentField.steep.isBelow(topLeft) {
		// case BLOCKING
		// Both lines intersect the square. This current field has ended.
		return activeFields.remove(currentField)

	} else if currentField.shallow.isAbove(bottomRight) {
		// case SHALLOW BUMP
		// The square intersects only the shallow line.
		addShallowBump(topLeft, currentField, shallowBumps)
		return checkField(currentField, activeFields)

	} else if currentField.steep.isBelow(topLeft) {
		// case STEEP BUMP
		// The square intersects only the steep line.
		addSteepBump(bottomRight, currentField, steepBumps)
		return checkField(currentField, activeFields)

	} else {
		// case BETWEEN
		// The square intersects neither line. We need to split into two fields.
		steeper := currentField
		shallower := activeFields.addBefore(currentField, *currentField)
		addSteepBump(bottomRight, shallower, steepBumps)
		checkField(shallower, activeFields)
		addShallowBump(topLeft, steeper, shallowBumps)
		return checkField(steeper, activeFields)
	}
}

func addShallowBump(point offset, currentField *field, shallowBumps *[]*bump) {
	// First, the far point of shallow is set to the new point.
	currentField.shallow.far = point
	// Second, we need to add the new bump to the shallow bump list for
	// future steep bump handling.
	newBump := bump{location: point, parent: currentField.shallowBump}
	*shallowBumps = append(*shallowBumps, &newBump)
	currentField.shallowBump = &newBump
	// Now we have too look through the list of steep bumps and see if
	// any of them are below the line.
	// If there are, we need to replace near point too.
	currentBump := currentField.steepBump
	for currentBump != nil {
		if currentField.shallow.isAbove(currentBump.location) {
			currentField.shallow.near = currentBump.location
		}
		currentBump = currentBump.parent
	}
}

func addSteepBump(point offset, currentField *field, steepBumps *[]*bump) {
	currentField.steep.far = point
	newBump := bump{location: point, parent: currentField.steepBump}
	*steepBumps = append(*steepBumps, &newBump)
	currentField.steepBump = &newBump

	// Now look through the list of shallow bumps and see if any of them
	// are below the line.
	currBump := currentField.shallowBump
	for currBump != nil {
		if currentField.steep.isBelow(currBump.location) {
			currentField.steep.near = currBump.location
		}
		currBump = currBump.parent
	}
}

func checkField(currentField *field, activeFields *fieldList) *field {
	result := currentField
	// If the two slopes are colinear, and if they pass through either
	// extremity, remove the field of view.
	if currentField.shallow.doesContain(currentField.steep.near) &&
		currentField.shallow.doesContain(currentField.steep.far) &&
		(currentField.shallow.doesContain(offset{0, 1}) || currentField.shallow.doesContain(offset{1, 0})) {
		result = currentField.next
		activeFields.remove(currentField)
	}
	return result
}


// TODO: delete the following
var imprint = 0
func immediateprint(msg string){
	console.PutString(fmt.Sprintf("%s - %d", msg, imprint), 0, 0)
	console.Flush_console()
	imprint++
	if imprint > 50000 {
		panic("Endless loop. I'd better crash.")
	}
}