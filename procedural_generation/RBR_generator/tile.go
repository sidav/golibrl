package RBR_generator

const (
	TWALL uint8 = iota
	TFLOOR
	TDOOR
	TUNKNOWN

	T_NOCHANGE 
)

type tile struct {
	tiletype byte
	roomId   int
	securityId int16 
}

func (t *tile) setProperties(ttype uint8, roomId int, secId int16) {
	if ttype != T_NOCHANGE {
		t.tiletype = ttype 
	}
	if roomId != -1 {
		t.roomId = roomId
	}
	if secId != -1 {
		t.securityId = secId 
	}
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
