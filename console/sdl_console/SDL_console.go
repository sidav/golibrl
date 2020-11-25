package sdl_console

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"io/ioutil"
	"os"
	"strings"
	"time"
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

	timeForMouseToBeHeld    = 75 * time.Millisecond
	timeForMouseToBeClicked = 25 * time.Millisecond
)

var (
	winTitle            = "Go-SDL2 Texture"
	chrW, chrH          int32
	termW, termH        int32 = 100, 40
	winWidth, winHeight int32
	FontPngFileName     string
	window              *sdl.Window
	renderer            *sdl.Renderer
	texture             *sdl.Texture
	fontImg             *sdl.Surface
	src, dst            sdl.Rect
	err                 error

	fgColor = []uint8{255, 255, 255}
	bgColor = []uint8{0, 0, 0}

	compatColorTable = map[int][]uint8{
		BLACK:        {0, 0, 0},
		DARK_RED:     {128, 0, 0},
		DARK_GREEN:   {0, 128, 0},
		DARK_YELLOW:  {128, 128, 0},
		DARK_BLUE:    {0, 0, 128},
		DARK_MAGENTA: {128, 0, 128},
		DARK_CYAN:    {0, 128, 128},
		BEIGE:        {128, 128, 96},
		DARK_GRAY:    {96, 96, 96},
		RED:          {255, 0, 0},
		GREEN:        {0, 255, 0},
		YELLOW:       {255, 255, 0},
		BLUE:         {0, 0, 255},
		MAGENTA:      {255, 0, 255},
		CYAN:         {0, 255, 255},
		WHITE:        {255, 255, 255},
	}

	wasResized     = false
	evCh           chan sdl.Event
	flushesCounter int

	mouseX, mouseY             int
	mouseVectorX, mouseVectorY int // for getting mouse coords changes
	lastMouseButton            string
	mouseClickedButton         string
	mouseMoved                 bool
	timeOfMousePress           time.Time

	// isShiftBeingHeld bool
)

func prepareFont() {
	files, err := ioutil.ReadDir("assets/")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if strings.Contains(f.Name(), ".png") {
			FontPngFileName = "assets/" + f.Name()
		}
	}

	fontImg, err = img.Load(FontPngFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load PNG: %s\n", err)
		return
	}
	chrW = fontImg.W / 16
	chrH = fontImg.H / 16
	winWidth, winHeight = termW*chrW, termH*chrH
	fontImg.SetColorKey(true, 0xff00ff)
}

func Init_console(title string) {

	prepareFont()

	window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN+sdl.WINDOW_RESIZABLE)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE) //sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return
	}

	texture, err = renderer.CreateTextureFromSurface(fontImg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		return
	}

	evCh = make(chan sdl.Event, 1)
	go startAsyncEventListener()

	renderer.Clear()
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.FillRect(&sdl.Rect{0, 0, int32(winWidth), int32(winHeight)})
	renderer.Copy(texture, &src, &dst)
	renderer.Present()
	sdl.Delay(0)
}

func Close_console() { //should be deferred!
	window.Destroy()
	renderer.Destroy()
	texture.Destroy()
	fontImg.Free()
}

func PurgeConsole() { // compatibility stub
	renderer.Clear()
}

func Clear_console() {
	//SetFgColorRGB(255, 255, 255)
	//SetBgColorRGB(0, 0, 0)
	//renderer.FillRect(&sdl.Rect{0, 0, winWidth, winHeight})
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.FillRect(&sdl.Rect{0, 0, int32(winWidth), int32(winHeight)})
}

func Flush_console() {
	renderer.Present()
	sdl.Delay(0)
}

func GetConsoleSize() (int, int) {
	return int(termW), int(termH)
}

func WasResized() bool { // stub for now
	return wasResized
}

func SetFgColorRGB(r, g, b uint8) {
	fgColor[0] = r
	fgColor[1] = g
	fgColor[2] = b
}

func SetBgColorRGB(r, g, b uint8) {
	bgColor[0] = r
	bgColor[1] = g
	bgColor[2] = b
}

func SetColor(fg int, bg int) {
	SetFgColor(fg)
	SetBgColor(bg)
}

func SetFgColor(fg int) {
	rgb := compatColorTable[fg]
	SetFgColorRGB(rgb[0], rgb[1], rgb[2])
}

func SetBgColor(bg int) {
	rgb := compatColorTable[bg]
	SetBgColorRGB(rgb[0], rgb[1], rgb[2])
}

func PutChar(c rune, x, y int) {
	code := int32(c)
	if code < 256 {
		row := code / 16 // chrH
		col := code % 16 // chrH
		// fmt.Printf("rune %s, code %d int32code %d, row/col %d, %d \n", string(c), int(c), code, row, col)
		src = sdl.Rect{chrW * col, chrH * row, chrW, chrH}
		dst = sdl.Rect{chrW * int32(x), chrH * int32(y), chrW, chrH}
		renderer.SetDrawColor(bgColor[0], bgColor[1], bgColor[2], 255)
		renderer.FillRect(&sdl.Rect{chrW * int32(x), chrH * int32(y), chrW, chrH})
		texture.SetColorMod(fgColor[0], fgColor[1], fgColor[2])
		renderer.Copy(texture, &src, &dst)
	}
}

func PutString(s string, x, y int) {
	length := len([]rune(s))
	for i := 0; i < length; i++ {
		PutChar([]rune(s)[i], x+i, y)
	}
}

func ReadKey() string {
	for {
		for len(evCh) == 0 { // wait here until an event is in the event queue
			time.Sleep(10 * time.Millisecond)
		}
		event := <-evCh
		switch t := event.(type) {
		case *sdl.WindowEvent:
			windowEventWork(t)
		case *sdl.KeyboardEvent:
			if t.State == 1 {
				keyString := sdl.GetScancodeName(t.Keysym.Scancode)

				// for compatibility...
				keyString = strings.Replace(keyString, "Keypad ", "", -1)

				if (t.Keysym.Mod&sdl.KMOD_SHIFT) != 1 && len(keyString) == 1 {
					return strings.ToLower(keyString)
				}
				return strings.ToUpper(keyString)
			}
		}
	}
}

func ReadKeyAsync() string { // also reads mouse events... TODO: think of if separate mouse events reader is needed.
	if len(evCh) == 0 {
		return "NOTHING"
	}
	evnt := <-evCh
	switch ev := evnt.(type) {
	case *sdl.KeyboardEvent:
		if ev.State == 1 {
			keyString := sdl.GetScancodeName(ev.Keysym.Scancode)

			// for compatibility...
			keyString = strings.Replace(keyString, "Keypad ", "", -1)
			if keyString == "Return" {
				return "ENTER" // TODO: reconsider?
			}

			if (ev.Keysym.Mod&sdl.KMOD_SHIFT) != 1 && len(keyString) == 1 {
				return strings.ToLower(keyString)
			}
			return strings.ToUpper(keyString)
		}
	case *sdl.MouseMotionEvent:
		mouseMoveWork(ev)
	case *sdl.MouseButtonEvent:
		mouseButtonWork(ev)
	case *sdl.WindowEvent:
		windowEventWork(ev)
	}
	return "NON-KEY"
}

func windowEventWork(wEvent *sdl.WindowEvent) {
	evnt := wEvent.Event
	if evnt == sdl.WINDOWEVENT_RESIZED || evnt == sdl.WINDOWEVENT_MOVED {
		win, _ := sdl.GetWindowFromID(wEvent.WindowID)
		winWidth, winHeight = win.GetSize()
		// winWidth, winHeight = wEvent.Data1, wEvent.Data2
		termW, termH = winWidth/chrW, winHeight/chrH
		window.SetSize(winWidth, winHeight)
		wasResized = true
	}
	Flush_console()
}

func mouseMoveWork(ev *sdl.MouseMotionEvent) {
	mx, my := int(ev.X/chrW), int(ev.Y/chrH)
	if mouseX != mx || mouseY != my {
		mouseVectorX = mx - mouseX
		mouseVectorY = my - mouseY
		mouseX, mouseY = mx, my
		mouseMoved = true
	}
}

func mouseButtonWork(ev *sdl.MouseButtonEvent) { // TODO: completely rewrite the method 
	// PrevMouseButton = mouseHeldButton
	var curMouseButton string
	if ev.Type == sdl.MOUSEBUTTONUP {
		curMouseButton = "NONE"
	} else {
		switch ev.Button {
		case sdl.BUTTON_LEFT:
			curMouseButton = "LEFT"
		case sdl.BUTTON_RIGHT:
			curMouseButton = "RIGHT"
		}
	}
	if curMouseButton != lastMouseButton {
		timeOfMousePress = time.Now()
	}
	timeSinceMousePress := time.Since(timeOfMousePress)
	// set click
	if curMouseButton == "NONE" {
		if timeSinceMousePress < timeForMouseToBeClicked {
			mouseClickedButton = lastMouseButton
		} else {
			mouseClickedButton = "NONE"
		}
	}
	lastMouseButton = curMouseButton
}

func GetMouseCoords() (int, int) {
	return mouseX, mouseY
}

func GetMouseHeldButton() string {
	timeSinceMousePress := time.Since(timeOfMousePress)
	if timeSinceMousePress >= timeForMouseToBeHeld {
		return lastMouseButton
	}
	return "NONE"
}

func GetMouseClickedButton() string {
	b := mouseClickedButton
	mouseClickedButton = "NONE"
	return b
}

func WasMouseMovedSinceLastEvent() bool {
	t := mouseMoved
	mouseMoved = false
	return t
}

func GetMouseMovementVector() (int, int) {
	return mouseVectorX, mouseVectorY
}

func startAsyncEventListener() {
	for {
		ev := sdl.WaitEvent()
		select {
		case evCh <- ev:
		default:
			sdl.Delay(0)
		}
	}
}

func PrintCharactersTable() {
	for x := 0; x < int(termW); x++ {
		for y := 0; y < int(termH); y++ {
			PutChar(rune(x+y*int(termW)), x, y)
		}
	}
	Flush_console()
	ReadKey()
}
