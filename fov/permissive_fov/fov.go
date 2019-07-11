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

func Fov(x, y, radius int) *[][]bool {
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
	const infinity = 1256// int(^uint32(0) >> 1)
	steepBumps := make([]bump, 0)
	shallowBumps := make([]bump, 0)
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
		current = activeFields.first
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
	return true
}

func visitSquare(dest offset, currentField *field, steepBumps *[]bump, shallowBumps *[]bump, activeFields *fieldList) *field {
	topLeft := offset{dest.x, dest.y + 1}
	bottomRight := offset{dest.x + 1, dest.y}
	for currentField != nil && currentField.steep.isBelowOrContains(bottomRight) {
		currentField = currentField.next
	}
	if currentField == nil || currentField.shallow.isAboveOrContains(topLeft) || !actIsBlocked(dest) {
		return currentField
	}
	if currentField.shallow.isAbove(bottomRight) && currentField.steep.isBelow(topLeft) {
		next := currentField.next
		activeFields.remove(currentField)
		return next
	} else if currentField.shallow.isAbove(bottomRight) {
		addShallowBump(topLeft, currentField, shallowBumps)
		return checkField(currentField, activeFields)
	} else if currentField.steep.isBelow(topLeft) {
		addSteepBump(bottomRight, currentField, steepBumps)
		return checkField(currentField, activeFields)
	} else {
		steeper := currentField
		shallower := activeFields.addBefore(currentField, *currentField)
		addSteepBump(bottomRight, shallower, steepBumps)
		checkField(shallower, activeFields)
		addShallowBump(topLeft, steeper, shallowBumps)
		return checkField(steeper, activeFields)
	}
}

func addShallowBump(point offset, currentField *field, shallowBumps *[]bump) {
	value := *currentField
	value.shallow.far = point
	value.shallowBump = &bump{value.shallowBump, point}
	*shallowBumps = append(*shallowBumps, *(value.shallowBump))

	currentBump := value.steepBump
	for currentBump != nil {
		if value.shallow.isAbove(currentBump.location) {
			value.shallow.near = currentBump.location
		}
		currentBump = currentBump.parent
	}
	*currentField = value
}

func addSteepBump(point offset, currentField *field, steepBumps *[]bump) {
	value := *currentField
	value.steep.far = point
	value.steepBump = &bump{location: point, parent: value.steepBump}
	*steepBumps = append(*steepBumps, *(value.steepBump))
	for currentBump := value.shallowBump; currentBump != nil; currentBump = currentBump.parent {
		if value.steep.isBelow(currentBump.location) {
			value.steep.near = currentBump.location
		}
	}
	*currentField = value
}

func checkField(currentField *field, activeFields *fieldList) *field {
	result := currentField
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
	return
	console.PutString(fmt.Sprintf("%s - %d", msg, imprint), 0, 0)
	console.Flush_console()
	imprint++
	if imprint > 50000 {
		panic("Endless loop. I'd better crash.")
	}
}