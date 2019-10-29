package RBR_generator

import (
	"bufio"
	"os"
	"strings"
	"github.com/sidav/golibrl/string_operations"
)

var vaults []*vault

type vault struct {
	strings []string
}

func (v *vault) getStrings() *[]string { // randomly rotates and/or mirrors the vault
	rotations := rnd.RandInRange(0, 3)
	vMirror := rnd.Rand(2) == 0
	hMirror := rnd.Rand(2) == 0
	result := string_operations.GetMirroredStringArray(&v.strings, vMirror, hMirror)
	for i := 1; i <= rotations; i++ {
		result = string_operations.GetRotatedStringArray(result)
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
