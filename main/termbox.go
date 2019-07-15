package main

import (
	"strconv"

	"github.com/limetext/backend/keys"
	"github.com/limetext/backend/log"
	"github.com/limetext/backend/render"
	"github.com/limetext/termbox-go"
)

var (
	lut = map[termbox.Key]keys.KeyPress{
		// Omission of these are intentional due to map collisions
		//		termbox.KeyCtrlTilde:      keys.KeyPress{Ctrl: true, Key: '~'},
		//		termbox.KeyCtrlBackslash:  keys.KeyPress{Ctrl: true, Key: '\\'},
		//		termbox.KeyCtrlSlash:      keys.KeyPress{Ctrl: true, Key: '/'},
		//		termbox.KeyCtrlUnderscore: keys.KeyPress{Ctrl: true, Key: '_'},
		//		termbox.KeyCtrlLsqBracket: keys.KeyPress{Ctrl: true, Key: '{'},
		//		termbox.KeyCtrlRsqBracket: keys.KeyPress{Ctrl: true, Key: '}'},
		// termbox.KeyCtrl3:
		// termbox.KeyCtrl8
		//		termbox.KeyCtrl2:      keys.KeyPress{Ctrl: true, Key: '2'},
		termbox.KeyCtrlSpace:  {Ctrl: true, Key: ' '},
		termbox.KeyCtrlA:      {Ctrl: true, Key: 'a'},
		termbox.KeyCtrlB:      {Ctrl: true, Key: 'b'},
		termbox.KeyCtrlC:      {Ctrl: true, Key: 'c'},
		termbox.KeyCtrlD:      {Ctrl: true, Key: 'd'},
		termbox.KeyCtrlE:      {Ctrl: true, Key: 'e'},
		termbox.KeyCtrlF:      {Ctrl: true, Key: 'f'},
		termbox.KeyCtrlG:      {Ctrl: true, Key: 'g'},
		termbox.KeyCtrlH:      {Ctrl: true, Key: 'h'},
		termbox.KeyCtrlJ:      {Ctrl: true, Key: 'j'},
		termbox.KeyCtrlK:      {Ctrl: true, Key: 'k'},
		termbox.KeyCtrlL:      {Ctrl: true, Key: 'l'},
		termbox.KeyCtrlN:      {Ctrl: true, Key: 'n'},
		termbox.KeyCtrlO:      {Ctrl: true, Key: 'o'},
		termbox.KeyCtrlP:      {Ctrl: true, Key: 'p'},
		termbox.KeyCtrlQ:      {Ctrl: true, Key: 'q'},
		termbox.KeyCtrlR:      {Ctrl: true, Key: 'r'},
		termbox.KeyCtrlS:      {Ctrl: true, Key: 's'},
		termbox.KeyCtrlT:      {Ctrl: true, Key: 't'},
		termbox.KeyCtrlU:      {Ctrl: true, Key: 'u'},
		termbox.KeyCtrlV:      {Ctrl: true, Key: 'v'},
		termbox.KeyCtrlW:      {Ctrl: true, Key: 'w'},
		termbox.KeyCtrlX:      {Ctrl: true, Key: 'x'},
		termbox.KeyCtrlY:      {Ctrl: true, Key: 'y'},
		termbox.KeyCtrlZ:      {Ctrl: true, Key: 'z'},
		termbox.KeyCtrl4:      {Ctrl: true, Key: '4'},
		termbox.KeyCtrl5:      {Ctrl: true, Key: '5'},
		termbox.KeyCtrl6:      {Ctrl: true, Key: '6'},
		termbox.KeyCtrl7:      {Ctrl: true, Key: '7'},
		termbox.KeyEnter:      {Key: keys.Enter},
		termbox.KeySpace:      {Key: ' '},
		termbox.KeyBackspace2: {Key: keys.Backspace},
		termbox.KeyArrowUp:    {Key: keys.Up},
		termbox.KeyArrowDown:  {Key: keys.Down},
		termbox.KeyArrowLeft:  {Key: keys.Left},
		termbox.KeyArrowRight: {Key: keys.Right},
		termbox.KeyDelete:     {Key: keys.Delete},
		termbox.KeyEsc:        {Key: keys.Escape},
		termbox.KeyPgup:       {Key: keys.PageUp},
		termbox.KeyPgdn:       {Key: keys.PageDown},
		termbox.KeyF1:         {Key: keys.F1},
		termbox.KeyF2:         {Key: keys.F2},
		termbox.KeyF3:         {Key: keys.F3},
		termbox.KeyF4:         {Key: keys.F4},
		termbox.KeyF5:         {Key: keys.F5},
		termbox.KeyF6:         {Key: keys.F6},
		termbox.KeyF7:         {Key: keys.F7},
		termbox.KeyF8:         {Key: keys.F8},
		termbox.KeyF9:         {Key: keys.F9},
		termbox.KeyF10:        {Key: keys.F10},
		termbox.KeyF11:        {Key: keys.F11},
		termbox.KeyF12:        {Key: keys.F12},
		termbox.KeyTab:        {Key: '\t'},
	}

	palLut    func(col render.Colour) termbox.Attribute
	defaultBg = termbox.ColorBlack
	defaultFg = termbox.ColorWhite
)

func addString(x, y int, s string, fg, bg termbox.Attribute) int {
	runes := []rune(s)
	addRunes(x, y, runes, fg, bg)
	x += len(runes)
	return x
}

func addRunes(x, y int, runes []rune, fg, bg termbox.Attribute) {
	for i, r := range runes {
		termbox.SetCell(x+i, y, r, fg, bg)
	}
}

func setColorMode() {
	var (
		mode256 bool
		pal     = make([]termbox.RGB, 0, 256)
	)

	if err := termbox.SetColorMode(termbox.ColorMode256); err != nil {
		log.Error("Unable to use 256 color mode: %s", err)
	} else {
		log.Debug("Using 256 color mode")
		mode256 = true
	}

	if !mode256 {
		pal = pal[:10] // Not correct, but whatever
		pal[termbox.ColorBlack] = termbox.RGB{R: 0, G: 0, B: 0}
		pal[termbox.ColorWhite] = termbox.RGB{R: 255, G: 255, B: 255}
		pal[termbox.ColorRed] = termbox.RGB{R: 255, G: 0, B: 0}
		pal[termbox.ColorGreen] = termbox.RGB{R: 0, G: 255, B: 0}
		pal[termbox.ColorBlue] = termbox.RGB{R: 0, G: 0, B: 255}
		pal[termbox.ColorMagenta] = termbox.RGB{R: 255, G: 0, B: 255}
		pal[termbox.ColorYellow] = termbox.RGB{R: 255, G: 255, B: 0}
		pal[termbox.ColorCyan] = termbox.RGB{R: 0, G: 255, B: 255}

		diff := func(i, j byte) int {
			v := int(i) - int(j)
			if v < 0 {
				return -v
			}
			return v
		}
		palLut = func(col render.Colour) termbox.Attribute {
			mindist := 10000000
			mini := 0
			for i, c := range pal {
				if dist := diff(c.R, col.R) + diff(c.G, col.G) + diff(c.B, col.B); dist < mindist {
					mindist = dist
					mini = i
				}
			}
			return termbox.Attribute(mini)
		}
	} else {
		palLut = func(col render.Colour) termbox.Attribute {
			tc := termbox.RGB{R: col.R, G: col.G, B: col.B}
			for i, c := range pal {
				if c == tc {
					return termbox.Attribute(i)
				}
			}
			l := len(pal)
			log.Debug("Adding colour: %d %+v %+v", l, col, tc)
			pal = append(pal, tc)
			termbox.SetColorPalette(pal)
			return termbox.Attribute(l)
		}
	}
}

func intToRunes(n int) (runes []rune) {
	lineStr := strconv.FormatInt(int64(n), 10)

	return []rune(lineStr)
}

func padLineRunes(line []rune, totalLineSize int) (padded []rune) {
	currentLineSize := len(line)
	if currentLineSize < totalLineSize {
		padding := (totalLineSize - currentLineSize)

		for i := 0; i < padding; i++ {
			padded = append(padded, ' ')
		}
	}

	padded = append(padded, line...)
	padded = append(padded, ' ')

	return
}

func renderLineNumber(line, x *int, y, lineNumberRenderSize int, fg, bg termbox.Attribute) {
	if *x == 0 {
		lineRunes := padLineRunes(intToRunes(*line), lineNumberRenderSize)

		for _, num := range lineRunes {
			termbox.SetCell(*x, y, num, fg, bg)
			*x++
		}

		*line++
	}

}

func getCaretStyle(style string, inverse bool) termbox.Attribute {
	caret_style := termbox.AttrUnderline

	if style == "block" {
		caret_style = termbox.AttrReverse
	}

	if inverse {
		if caret_style == termbox.AttrReverse {
			caret_style = termbox.AttrUnderline
		} else {
			caret_style = termbox.AttrReverse
		}
	}

	return caret_style
}
