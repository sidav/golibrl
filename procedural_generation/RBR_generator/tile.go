package RBR_generator

import "strconv"

const (
	TWALL uint8 = iota
	TFLOOR
	TDOOR
	TNEXTLEVELSTAIR
	TPREVLEVELSTAIR
	TUNKNOWN

	T_NOCHANGE
)

type tile struct {
	TileType byte
	roomId   int
	SecArea  int16
}

func (t *tile) setProperties(ttype uint8, roomId int, secId int16) {
	if ttype != T_NOCHANGE {
		t.TileType = ttype
	}
	if roomId != -1 {
		t.roomId = roomId
	}
	if secId != -1 {
		t.SecArea = secId
	}
}

func (t *tile) toRune() rune {
	switch t.TileType {
	case TFLOOR:
		return '.'
	case TWALL:
		return '#'
	case TNEXTLEVELSTAIR:
		return '>'
	case TPREVLEVELSTAIR:
		return '<'
	case TDOOR:
		if t.SecArea == 0 {
			return '+'
		} else {
			return rune(strconv.Itoa(int(t.SecArea))[0]) //'\\'
		}
	}
	return '?'
}
