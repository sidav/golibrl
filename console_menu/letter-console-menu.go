package console_menu 

import (
	cw "github.com/sidav/golibrl/console"
	// "strings"
	"sort"
)

func GetSortedKeysListForMap(m *map[rune]string) *[]rune {
	list := make([]rune, 0)
	for k, _ := range *m {
		list = append(list, k)
	}
	sort.Sort()
	return &list 
}

func ShowSingleChoiceLetterMenu(title, subheading string, menuItems *map[rune]string, allowcursor bool) rune { //returns the key of selected line or -1 if nothing was selected.
	// val := lines
	cursor := 0
	sortedKeys := GetSortedKeysListForMap(menuItems)
	for {
		cw.Clear_console()
		drawTitle(title)
		cw.SetFgColor(cw.BEIGE)
		cw.PutString(subheading, 0, 1)
		for i, k := range *sortedKeys {
			if cursor == i {
				cw.SetColor(cw.BLACK, TEXT_COLOR)
			} else {
				cw.SetColor(TEXT_COLOR, cw.BLACK)
			}
			cw.PutString(" "+string(k) + " - " +v+ " ", 1, 2+i)
			cw.SetBgColor(cw.BLACK)
		}
		cw.Flush_console()
		key := cw.ReadKey()
		switch key {
		case "2":
			cursor++
			if cursor == len(val) {
				cursor = 0
			}
		case "8":
			cursor--
			if cursor < 0 {
				cursor = len(val) - 1
			}
		case "ENTER":
			return cursor
		case "ESCAPE":
			return -1
		}
	}
}

