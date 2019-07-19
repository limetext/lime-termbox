package main

import (
	"strconv"

	"github.com/limetext/backend/keys"
	"github.com/limetext/backend/render"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/nsf/termbox-go"
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

	// xterm 256 colors
	// https://jonasjacek.github.io/colors/
	palette = []string{
		"#000000",
		"#800000",
		"#008000",
		"#808000",
		"#000080",
		"#800080",
		"#008080",
		"#c0c0c0",
		"#808080",
		"#ff0000",
		"#00ff00",
		"#ffff00",
		"#0000ff",
		"#ff00ff",
		"#00ffff",
		"#ffffff",
		"#000000",
		"#00005f",
		"#000087",
		"#0000af",
		"#0000d7",
		"#0000ff",
		"#005f00",
		"#005f5f",
		"#005f87",
		"#005faf",
		"#005fd7",
		"#005fff",
		"#008700",
		"#00875f",
		"#008787",
		"#0087af",
		"#0087d7",
		"#0087ff",
		"#00af00",
		"#00af5f",
		"#00af87",
		"#00afaf",
		"#00afd7",
		"#00afff",
		"#00d700",
		"#00d75f",
		"#00d787",
		"#00d7af",
		"#00d7d7",
		"#00d7ff",
		"#00ff00",
		"#00ff5f",
		"#00ff87",
		"#00ffaf",
		"#00ffd7",
		"#00ffff",
		"#5f0000",
		"#5f005f",
		"#5f0087",
		"#5f00af",
		"#5f00d7",
		"#5f00ff",
		"#5f5f00",
		"#5f5f5f",
		"#5f5f87",
		"#5f5faf",
		"#5f5fd7",
		"#5f5fff",
		"#5f8700",
		"#5f875f",
		"#5f8787",
		"#5f87af",
		"#5f87d7",
		"#5f87ff",
		"#5faf00",
		"#5faf5f",
		"#5faf87",
		"#5fafaf",
		"#5fafd7",
		"#5fafff",
		"#5fd700",
		"#5fd75f",
		"#5fd787",
		"#5fd7af",
		"#5fd7d7",
		"#5fd7ff",
		"#5fff00",
		"#5fff5f",
		"#5fff87",
		"#5fffaf",
		"#5fffd7",
		"#5fffff",
		"#870000",
		"#87005f",
		"#870087",
		"#8700af",
		"#8700d7",
		"#8700ff",
		"#875f00",
		"#875f5f",
		"#875f87",
		"#875faf",
		"#875fd7",
		"#875fff",
		"#878700",
		"#87875f",
		"#878787",
		"#8787af",
		"#8787d7",
		"#8787ff",
		"#87af00",
		"#87af5f",
		"#87af87",
		"#87afaf",
		"#87afd7",
		"#87afff",
		"#87d700",
		"#87d75f",
		"#87d787",
		"#87d7af",
		"#87d7d7",
		"#87d7ff",
		"#87ff00",
		"#87ff5f",
		"#87ff87",
		"#87ffaf",
		"#87ffd7",
		"#87ffff",
		"#af0000",
		"#af005f",
		"#af0087",
		"#af00af",
		"#af00d7",
		"#af00ff",
		"#af5f00",
		"#af5f5f",
		"#af5f87",
		"#af5faf",
		"#af5fd7",
		"#af5fff",
		"#af8700",
		"#af875f",
		"#af8787",
		"#af87af",
		"#af87d7",
		"#af87ff",
		"#afaf00",
		"#afaf5f",
		"#afaf87",
		"#afafaf",
		"#afafd7",
		"#afafff",
		"#afd700",
		"#afd75f",
		"#afd787",
		"#afd7af",
		"#afd7d7",
		"#afd7ff",
		"#afff00",
		"#afff5f",
		"#afff87",
		"#afffaf",
		"#afffd7",
		"#afffff",
		"#d70000",
		"#d7005f",
		"#d70087",
		"#d700af",
		"#d700d7",
		"#d700ff",
		"#d75f00",
		"#d75f5f",
		"#d75f87",
		"#d75faf",
		"#d75fd7",
		"#d75fff",
		"#d78700",
		"#d7875f",
		"#d78787",
		"#d787af",
		"#d787d7",
		"#d787ff",
		"#d7af00",
		"#d7af5f",
		"#d7af87",
		"#d7afaf",
		"#d7afd7",
		"#d7afff",
		"#d7d700",
		"#d7d75f",
		"#d7d787",
		"#d7d7af",
		"#d7d7d7",
		"#d7d7ff",
		"#d7ff00",
		"#d7ff5f",
		"#d7ff87",
		"#d7ffaf",
		"#d7ffd7",
		"#d7ffff",
		"#ff0000",
		"#ff005f",
		"#ff0087",
		"#ff00af",
		"#ff00d7",
		"#ff00ff",
		"#ff5f00",
		"#ff5f5f",
		"#ff5f87",
		"#ff5faf",
		"#ff5fd7",
		"#ff5fff",
		"#ff8700",
		"#ff875f",
		"#ff8787",
		"#ff87af",
		"#ff87d7",
		"#ff87ff",
		"#ffaf00",
		"#ffaf5f",
		"#ffaf87",
		"#ffafaf",
		"#ffafd7",
		"#ffafff",
		"#ffd700",
		"#ffd75f",
		"#ffd787",
		"#ffd7af",
		"#ffd7d7",
		"#ffd7ff",
		"#ffff00",
		"#ffff5f",
		"#ffff87",
		"#ffffaf",
		"#ffffd7",
		"#ffffff",
		"#080808",
		"#121212",
		"#1c1c1c",
		"#262626",
		"#303030",
		"#3a3a3a",
		"#444444",
		"#4e4e4e",
		"#585858",
		"#626262",
		"#6c6c6c",
		"#767676",
		"#808080",
		"#8a8a8a",
		"#949494",
		"#9e9e9e",
		"#a8a8a8",
		"#b2b2b2",
		"#bcbcbc",
		"#c6c6c6",
		"#d0d0d0",
		"#dadada",
		"#e4e4e4",
		"#eeeeee",
	}
	colorMap = map[string]termbox.Attribute{}

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
	termbox.SetOutputMode(termbox.Output256)
}

func color256(col render.Colour) termbox.Attribute {
	if attr, ok := colorMap[col.String()]; ok {
		return attr
	}

	c1 := colorful.Color{
		R: float64(col.R) / 255.0,
		G: float64(col.G) / 255.0,
		B: float64(col.B) / 255.0,
	}
	dist := float64(1)
	found := 1
	for i, hex := range palette {
		c2, _ := colorful.Hex(hex)
		nd := c1.DistanceLab(c2)
		if nd < dist {
			dist = nd
			found = i + 1
		}
	}

	attr := termbox.Attribute(found)
	colorMap[col.String()] = attr

	return attr
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
