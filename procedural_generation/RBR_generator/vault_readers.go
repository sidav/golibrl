package RBR_generator

import (
	"bufio"
	"os"
	"strings"
)

func (r *RBR) readRoomVaultsFromFile(path string) {
	r.roomvaults = make([]*vault, 0)
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
				r.roomvaults = append(r.roomvaults, &vault{strings: vaultLines})
				vaultLines = make([]string, 0)
			}
		} else {
			vaultLines = append(vaultLines, line)
		}
	}
	if len(vaultLines) > 0 {
		r.roomvaults = append(r.roomvaults, &vault{strings: vaultLines})
	}
}

func (r *RBR) readVaultsFromFile(path string) {
	r.vaults = make([]*vault, 0)
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
				r.vaults = append(r.vaults, &vault{strings: vaultLines})
				vaultLines = make([]string, 0)
			}
		} else {
			vaultLines = append(vaultLines, line)
		}
	}
	if len(vaultLines) > 0 {
		r.vaults = append(r.vaults, &vault{strings: vaultLines})
	}
}

func (r *RBR) getRandomVault() *vault {
	return r.vaults[rnd.Rand(len(r.vaults))]
}

func (r *RBR) getRandomRoomvault() *vault {
	return r.roomvaults[rnd.Rand(len(r.roomvaults))]
}
