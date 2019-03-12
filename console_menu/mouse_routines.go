package console_menu

import cw "github.com/sidav/goLibRL/console"

func isMouseInMenuBounds(mx, my, mw, mh int) bool {
	mousex, mousey := cw.GetMouseCoords()
	return mousex >= mx && mousex < mx+mw && mousey > my && mousey < my + mh
}
