package RBR_generator

import (
	"bufio"
	"os"
	"strings"
)

var vaults []*vault

type vault struct {
	strings []string
}

func reverseString(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

func getRotatedStringArray(arr *[]string) *[]string { // rotates 90 degrees clockwise
	newArr := make([]string, 0)
	for i := 0; i < len((*arr)[0]); i++ {
		str := ""
		for j := 0; j < len(*arr); j++ {
			str += string((*arr)[j][i])
		}
		newArr = append(newArr, str)
	}
	return &newArr
}

func getMirroredStringArray(arr *[]string, v, h bool) *[]string {
	newArr := make([]string, 0)
	if v && h {
		for i := len(*arr) - 1; i >= 0; i-- {
			newArr = append(newArr, reverseString((*arr)[i]))
		}
		return &newArr
	}
	if v {
		for i := len(*arr) - 1; i >= 0; i-- {
			newArr = append(newArr, (*arr)[i])
		}
		return &newArr
	}
	if h {
		for i := 0; i < len(*arr); i++ {
			newArr = append(newArr, reverseString((*arr)[i]))
		}
		return &newArr
	}
	return arr // nil
}

func (v *vault) getStrings() *[]string { // randomly rotates and/or mirrors the vault 
	rotations := rnd.RandInRange(0, 3)
	vMirror := rnd.Rand(2) == 0
	hMirror := rnd.Rand(2) == 0
	result := getMirroredStringArray(&v.strings, vMirror, hMirror)
	for i := 1; i <= rotations; i++ {
		result = getRotatedStringArray(result)
	}
	return result 
}

func readVaultsFromFile(path string) {
	vaults = make([]*vault, 0)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	vaultLines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.Contains(line, "//") {
			if len(vaultLines) > 0 {
				vaults = append(vaults, &vault{strings: vaultLines})
				vaultLines = make([]string, 0)
			}
		} else {
			vaultLines = append(vaultLines, line)
		}
	}
	if len(vaultLines) > 0 {
		vaults = append(vaults, &vault{strings: vaultLines})
	}
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
