package RBR_generator

import (
	"github.com/sidav/golibrl/string_operations"
)

type vault struct {
	strings []string
}

func (v *vault) getStrings() *[]string { // randomly rotates and/or mirrors the vault
	rotations := rnd.RandInRange(0, 3)
	vMirror := rnd.OneChanceFrom(2)
	hMirror := rnd.OneChanceFrom(2)
	result := string_operations.GetMirroredStringArray(&v.strings, vMirror, hMirror)
	for i := 1; i <= rotations; i++ {
		result = string_operations.GetRotatedStringArray(result)
	}
	return result
}

func (v *vault) isOfSize(w, h int) bool {
	vh, vw := len(v.strings), len(v.strings[0])
	if vw == w && vh == h { 
		return true 
	}
	if vw == h && vh == w { // will be fit if rotated 
		return true 
	}
	return false 
}

func (v *vault) getStringsIfFitInSize(w, h int) *[]string {
	vh, vw := len(v.strings), len(v.strings[0])
	if vw == w && vh == h && vw == vh { // square vault, fits
		return v.getStrings()
	}
	if vw == w && vh == h { // only mirror 
		vMirror := rnd.OneChanceFrom(2)
		hMirror := rnd.OneChanceFrom(2)
		result := string_operations.GetMirroredStringArray(&v.strings, vMirror, hMirror)
		return result
	}
	if vw == h && vh == w { // will be fit if rotated 
		vMirror := rnd.OneChanceFrom(2)
		hMirror := rnd.OneChanceFrom(2)
		result := string_operations.GetMirroredStringArray(&v.strings, vMirror, hMirror)
		result = string_operations.GetRotatedStringArray(result)
		return result
	}
	return nil 
}

func vaultSymbolToTileType(symbol rune) byte {
	switch symbol {
	case '#':
		return TWALL
	case '+':
		return TDOOR
	case '.':
		return TFLOOR
	default:
		return TUNKNOWN
	}
}
