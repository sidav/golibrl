package permissive_fov

type offset struct {
	x, y int
}

type bump struct {
	parent *bump
	location offset
}

type field struct {
	steepBump, shallowBump *bump
	steep, shallow line
	next *field // it's like a linked list
	prev *field
}

type fieldList struct {
	first *field
	last *field
}

func (fl *fieldList) addToEnd(f *field) {
	f.prev = fl.last
	fl.last = f
}

func (fl *fieldList) addToBeginning(f *field) {
	f.next = fl.first
	fl.first = f
}

func (fl *fieldList) remove(f *field) {
	curr := fl.first
	for curr != nil {
		if curr == f {
			prev := curr.prev
			next := curr.next
			if prev == nil {
				fl.first = next
			}
			if next == nil {
				fl.last = prev 
			}
			return
		}
		curr = curr.next
	}
}

type line struct {
	near, far offset
}

func (l *line) relativeSlope(point offset) int {
	return (l.far.y - l.near.y)*(l.far.x - point.x) - (l.far.y - point.y)*(l.far.x - l.near.x)
}

func (l *line) isBelow(point offset) bool {
	return l.relativeSlope(point) > 0
}

func (l *line) isBelowOrContains(point offset) bool {
	return l.relativeSlope(point) >= 0
}

func (l *line) isAbove(point offset) bool {
	return l.relativeSlope(point) < 0
}

func (l *line) isAboveOrContains(point offset) bool {
	return l.relativeSlope(point) <= 0
}

func (l *line) doesContain(point offset) bool {
	return l.relativeSlope(point) == 0
}
