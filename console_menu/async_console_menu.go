package console_menu

import (
	cw "github.com/sidav/goLibRL/console"
	"strings"
)

// returns index of element. Is not looped, does not flush the console and is purposed for alongside-usage.
func DrawSidebarMouseOnlyAsyncMenu(title string, titleColor, mx, my, mw int, items []string) int{
	drawSidebarMenuTitle(title, titleColor, mx, my, mw)
	mh := len(items)+1 // +1 is for the title

	cursorIndex := -1
	_, mousey := cw.GetMouseCoords()
	if isMouseInMenuBounds(mx, my, mw, mh) {
		cursorIndex = mousey - my - 1
	}

	for y := 1; y < mh; y++ {
		cw.PutString(strings.Repeat(" ", mw), mx, y+my) // clear menu screen space
	}
	for i := 0; i < len(items); i++ {
		str := items[i]
		if i == cursorIndex {
			cw.SetBgColor(cw.BEIGE)
			cw.SetFgColor(cw.BLACK)
			// str = "->"+str
		} else {
			cw.SetBgColor(cw.BLACK)
			cw.SetFgColor(cw.BEIGE)
		}
		// str += strings.Repeat(" ", mw - len(str)) // fill the whole menu width
		cw.PutString(str, mx, my+i+1)
		cw.SetBgColor(cw.BLACK)
		cw.SetFgColor(cw.BEIGE)
	}
	cw.SetBgColor(cw.BLACK)
	return cursorIndex
}
