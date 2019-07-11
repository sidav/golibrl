package permissive_fov

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
	const infinity = int(^uint32(0) >> 1)
	steepBumps := make([]bump, 0)
	shallowBumps := make([]bump, 0)
	activeFields := make([]field, 0)
	activeFields = append(activeFields, field{
		steep:   line{offset{1, 0}, offset{0, infinity}},
		shallow: line{offset{0, 1}, offset{infinity, 0}},
	})
	dest := offset{}
	actIsBlocked(dest)
	for i := 1; i < infinity && len(activeFields) > 0; i++ {
		current := &activeFields[0]
		for j := 0; j <= i; j++ {
			dest.x = i - j
			dest.y = j
			current = visitSquare(dest, current, steepBumps, shallowBumps, activeFields)
		}
	}
}

func actIsBlocked(pos offset) bool {
	if rangeLimit >= 0 && getDistance(max(pos.x, pos.y), min(pos.x, pos.y)) > rangeLimit {
		return true
	}
	x := pos.x*quadrant.x + source.x
	y := pos.y*quadrant.y + source.y
	(*visible)[x][y] = true
	return (*opaque)[x][y]
}

func visitSquare(dest offset, currentField *field, steepBumps []bump, shallowBumps []bump, activeFields []field) *field {
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
		// activeFields = remove(activeFields)
		return next
	} else if currentField.shallow.isAbove(bottomRight) {
		addShallowBump(topLeft, currentField, shallowBumps)
		return checkField(currentField, activeFields)
	} else if currentField.steep.isBelow(topLeft) {
		addSteepBump(bottomRight, currentField, steepBumps)
		return checkField(currentField, activeFields)
	} else {
		steeper := currentField
		// shallower := activeFields.AddBefore(currentField, currentField)
		// addSteepBump(bottomRight, shallower, steepBumps)
		// checkField(shallower, activeFields)
		addShallowBump(topLeft, steeper, shallowBumps)
		return checkField(steeper, activeFields)
	}
	return nil
}

func addShallowBump(point offset, currentField *field, shallowBumps []bump) {
	value := currentField
	value.shallow.far = point
	value.shallowBump = &bump{value.shallowBump, point}
	shallowBumps = append(shallowBumps, *(value.shallowBump))

	currentBump := value.steepBump
	for currentBump != nil {
		if value.shallow.isAbove(currentBump.location) {
			value.shallow.near = currentBump.location
		}
		currentBump = currentBump.parent
	}
	currentField = value
}

func addSteepBump(point offset, currentField *field, steepBumps []bump) {
	value := currentField
	value.steep.far = point
	value.steepBump = &bump{location: point, parent: value.steepBump}
	steepBumps = append(steepBumps, *value.steepBump)
	for currentBump := value.shallowBump; currentBump != nil; currentBump = currentBump.parent {
		if value.steep.isBelow(currentBump.location) {
			value.steep.near = currentBump.location
		}
	}
	currentField = value
}

func checkField(currentField *field, activeFields []field) *field {
	result := currentField
	if currentField.shallow.doesContain(currentField.steep.near) &&
		currentField.shallow.doesContain(currentField.steep.far) &&
		(currentField.shallow.doesContain(offset{0, 1}) || currentField.shallow.doesContain(offset{1, 0})) {
		result = currentField.next
		// activeFields.Remove(currentField)
	}
	return result
}
