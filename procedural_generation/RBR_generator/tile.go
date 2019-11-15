package RBR_generator
import "strconv"

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
	secArea int16 
}

func (t *tile) setProperties(ttype uint8, roomId int, secId int16) {
	if ttype != T_NOCHANGE {
		t.tiletype = ttype 
	}
	if roomId != -1 {
		t.roomId = roomId
	}
	if secId != -1 {
		t.secArea = secId 
	}
}

func (t *tile) toRune() rune {
	switch t.tiletype {
	case TFLOOR:
		return '.'
	case TWALL:
		return '#'
	case TDOOR:
		if t.secArea == 0 {
		return '+'
		} else {
			return rune(strconv.Itoa(int(t.secArea))[0]) //'\\'
		}
	}
	return '?'
}
