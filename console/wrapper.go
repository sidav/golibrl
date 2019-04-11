package console

import (
	"github.com/sidav/golibrl/console/sdl_console"
	"github.com/sidav/golibrl/console/tcell_console"
)

const (
	TCellRenderer = iota
	SDLRenderer   = iota
)

const ( // for the great compatibility with default console color codes
	BLACK        = 0
	DARK_RED     = 1
	DARK_GREEN   = 2
	DARK_YELLOW  = 3
	DARK_BLUE    = 4
	DARK_MAGENTA = 5
	DARK_CYAN    = 6
	BEIGE        = 7
	DARK_GRAY    = 8
	RED          = 9
	GREEN        = 10
	YELLOW       = 11
	BLUE         = 12
	MAGENTA      = 13
	CYAN         = 14
	WHITE        = 15
)

var (
	selectedRenderer                    = TCellRenderer
	flushesCounter int
	// isShiftBeingHeld bool
)

func Init_console(title string, preferableRenderer int) {
	selectedRenderer = preferableRenderer
	if selectedRenderer == SDLRenderer {
		sdl_console.Init_console(title)
		return
	}
	if selectedRenderer == TCellRenderer {
		tcell_console.Init_console()
	}
}

func Close_console() { //should be deferred!
	if selectedRenderer == SDLRenderer {
		sdl_console.Close_console()
		return
	}
	if selectedRenderer == TCellRenderer {
		tcell_console.Close_console()
	}
}

func PurgeConsole() { // compatibility stub
	if selectedRenderer == TCellRenderer {
		tcell_console.PurgeConsole()
	}
}

func Clear_console() {
	if selectedRenderer == SDLRenderer {
		sdl_console.Clear_console()
		return
	}
	if selectedRenderer == TCellRenderer {
		tcell_console.Clear_console()
	}
}

func Flush_console() {
	flushesCounter += 1
	if selectedRenderer == SDLRenderer {
		sdl_console.Flush_console()
		return
	}
	if selectedRenderer == TCellRenderer {
		tcell_console.Flush_console()
	}
}

func GetConsoleSize() (int, int) {
	if selectedRenderer == SDLRenderer {
		return sdl_console.GetConsoleSize()
	}
	if selectedRenderer == TCellRenderer {
		return tcell_console.GetConsoleSize()
	}
	return -1, -1
}

func WasResized() bool { // stub for now
	if selectedRenderer == SDLRenderer {
		return sdl_console.WasResized()
	}
	if selectedRenderer == TCellRenderer {
		return tcell_console.WasResized()
	}
	return false
}

func SetFgColorRGB(r, g, b uint8) {
	if selectedRenderer == SDLRenderer {
		sdl_console.SetFgColorRGB(r, g, b)
		return
	}
	if selectedRenderer == TCellRenderer {
		// TODO: rgb to 4 bits conversion lol
	}
}

func SetBgColorRGB(r, g, b uint8) {
	if selectedRenderer == SDLRenderer {
		sdl_console.SetBgColorRGB(r, g, b)
		return
	}
	if selectedRenderer == TCellRenderer {
		// TODO: rgb to 4 bits conversion lol
	}
}

func SetColor(fg int, bg int) {
	if selectedRenderer == SDLRenderer {
		sdl_console.SetColor(fg, bg)
		return
	}
	if selectedRenderer == TCellRenderer {
		tcell_console.SetColor(fg, bg)
	}
}

func SetFgColor(fg int) {
	if selectedRenderer == SDLRenderer {
		sdl_console.SetFgColor(fg)
		return
	}
	if selectedRenderer == TCellRenderer {
		tcell_console.SetFgColor(fg)
	}
}

func SetBgColor(bg int) {
	if selectedRenderer == SDLRenderer {
		sdl_console.SetBgColor(bg)
		return
	}
	if selectedRenderer == TCellRenderer {
		tcell_console.SetBgColor(bg)
	}
}

func PutChar(c rune, x, y int) {
	if selectedRenderer == SDLRenderer {
		sdl_console.PutChar(c, x, y)
		return
	}
	if selectedRenderer == TCellRenderer {
		tcell_console.PutChar(c, x, y)
	}
}

func PutString(s string, x, y int) {
	if selectedRenderer == SDLRenderer {
		sdl_console.PutString(s, x, y)
		return
	}
	if selectedRenderer == TCellRenderer {
		tcell_console.PutString(s, x, y)
	}
}

func ReadKey() string {
	if selectedRenderer == SDLRenderer {
		return sdl_console.ReadKey()
	}
	if selectedRenderer == TCellRenderer {
		return tcell_console.ReadKey()
	}
	return "WRAPPER ERROR"
}

func ReadKeyAsync() string { // also reads mouse events... TODO: think of if separate mouse events reader is needed.
	if selectedRenderer == SDLRenderer {
		return sdl_console.ReadKeyAsync()
	}
	if selectedRenderer == TCellRenderer {
		return tcell_console.ReadKeyAsync()
	}
	return "WRAPPER ERROR"
}

func GetMouseCoords() (int, int) {
	if selectedRenderer == SDLRenderer {
		return sdl_console.GetMouseCoords()
	}
	if selectedRenderer == TCellRenderer {
		return tcell_console.GetMouseCoords()
	}
	return -1, -1
}

func GetMouseHeldButton() string {
	if selectedRenderer == SDLRenderer {
		return sdl_console.GetMouseButton()
	}
	if selectedRenderer == TCellRenderer {
		return tcell_console.GetMouseHeldButton()
	}
	return "WRAPPER ERROR"
}

func GetMouseClickedButton() string {
	if selectedRenderer == SDLRenderer {
		return sdl_console.GetMouseButton()
	}
	if selectedRenderer == TCellRenderer {
		return tcell_console.GetMouseClickedButton()
	}
	return "WRAPPER ERROR"
}

func WasMouseMovedSinceLastEvent() bool {
	if selectedRenderer == SDLRenderer {
		return sdl_console.WasMouseMovedSinceLastEvent()
	}
	if selectedRenderer == TCellRenderer {
		return tcell_console.WasMouseMovedSinceLastEvent()
	}
	return false
}

func GetMouseMovementVector() (int, int) {
	if selectedRenderer == SDLRenderer {
		return sdl_console.GetMouseMovementVector()
	}
	if selectedRenderer == TCellRenderer {
		return tcell_console.GetMouseMovementVector()
	}
	return 0, 0
}


func GetNumberOfRecentFlushes() int { // may be useful for searching rendering overkills and something
	t := flushesCounter
	flushesCounter = 0
	return t
}

func PrintCharactersTable() {
	if selectedRenderer == SDLRenderer {
		sdl_console.PrintCharactersTable()
	}
	if selectedRenderer == TCellRenderer {
		tcell_console.PrintCharactersTable()
	}
}
