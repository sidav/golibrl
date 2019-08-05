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
	size int
}

func (fl *fieldList) addToEnd(f *field) {
	f.prev = fl.last
	if fl.last != nil {
		fl.last.next = f
	} else {
		fl.first = f
	}
	fl.last = f
	fl.size++
}

func (fldlist *fieldList) addBefore(f1 *field, f2 field) *field {
	curr := fldlist.first
	for curr != nil {
		if curr == f1 {
			prev := f1.prev
			if prev == nil {
				fldlist.first = &f2
			} else {
				prev.next = &f2
			}
			f2.next = f1
			f2.prev = prev
			curr.prev = &f2
			fldlist.size++
			if f1 == &f2 {
				panic("AddBefore has seen some strange shit.")
			}
			return &f2
		}
		curr = curr.next
	}
	panic("AddBefore has crashed.")
	return nil
}

func (fl *fieldList) remove(f *field) *field {
	curr := fl.first
	for curr != nil {
		if curr == f {
			prev := curr.prev
			next := curr.next
			if prev == nil {
				fl.first = next
			} else {
				prev.next = next
			}
			if next == nil {
				fl.last = prev
			} else {
				next.prev = prev
			}
			fl.size--
			return f.next
		}
		curr = curr.next
	}
	return nil
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
