package RBR_generator

const (
	TWALL = iota
	TFLOOR
	TDOOR
)

type tile struct {
	tiletype byte
	roomId   int
}

func (t *tile) toRune() rune {
	switch t.tiletype {
	case TFLOOR:
		return '.'
	case TWALL:
		return '#'
	case TDOOR:
		return '+'
	}
	return '?'
}
