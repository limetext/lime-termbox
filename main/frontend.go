package main

import (
	"flag"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/limetext/backend"
	"github.com/limetext/backend/keys"
	"github.com/limetext/backend/log"
	"github.com/limetext/backend/render"
	. "github.com/limetext/text"
	"github.com/limetext/util"
	"github.com/nsf/termbox-go"
)

type (
	tbfe struct {
		layout         map[*backend.View]layout
		window_layout  layout
		status_message string
		dorender       chan bool
		shutdown       chan bool
		lock           sync.Mutex
		editor         *backend.Editor
		console        *backend.View
		currentView    *backend.View
		currentWindow  *backend.Window
	}

	layout struct {
		x, y          int
		width, height int
		visible       Region
		lastUpdate    int
	}

	tbfeBufferDeltaObserver struct {
		t    *tbfe
		view *backend.View
	}
)

const (
	render_chan_len = 2
	statusbarHeight = 1
)

var (
	blink bool
)

func createFrontend() *tbfe {
	var t tbfe
	t.dorender = make(chan bool, render_chan_len)
	t.shutdown = make(chan bool, 2)
	t.layout = make(map[*backend.View]layout)

	t.editor = t.setupEditor()
	t.console = t.editor.Console()
	t.currentWindow = t.editor.NewWindow()

	// Assuming that all extra arguments are files
	if files := flag.Args(); len(files) > 0 {
		for _, file := range files {
			t.currentView = createNewView(file, t.currentWindow)
		}
	} else {
		t.currentView = t.currentWindow.NewFile()
	}

	t.editor.AddPackagesPath("../packages")

	t.editor.SetFrontend(&t)
	t.editor.LogInput(false)
	t.editor.LogCommands(false)

	w, h := termbox.Size()
	t.handleResize(h, w, true)

	t.console.AddObserver(&t)
	t.setupCallbacks(t.currentView)

	setColorMode()
	setSchemeSettings(t.editor)

	return &t
}

func (t *tbfe) renderView(v *backend.View, lay layout) {
	p := util.Prof.Enter("render")
	defer p.Exit()

	sx, sy, w, h := lay.x, lay.y, lay.width, lay.height
	vr := lay.visible
	runes := v.Substr(vr)
	x, y := sx, sy
	ex, ey := sx+w, sy+h

	style, _ := v.Settings().Get("caret_style", "underline").(string)
	inverse, _ := v.Settings().Get("inverse_caret_state", false).(bool)
	caretStyle := getCaretStyle(style, inverse)
	oldCaretStyle := caretStyle
	caretBlink, _ := v.Settings().Get("caret_blink", true).(bool)
	if caretBlink && blink {
		caretStyle = 0
	}

	tabSize := 4
	if i, ok := v.Settings().Get("tab_size", tabSize).(int); ok {
		tabSize = i
	}

	recipe := v.Transform(vr).Transcribe()
	fg, bg := defaultFg, defaultBg

	sel := v.Sel()

	lineNumbers, _ := v.Settings().Get("line_numbers", true).(bool)
	eofline, _ := v.RowCol(v.Size())
	lineNumberRenderSize := len(intToRunes(eofline))
	line, _ := v.RowCol(vr.Begin())
	line += 1

	for i, r := range runes {
		fg, bg = defaultFg, defaultBg

		if lineNumbers {
			renderLineNumber(&line, &x, y, lineNumberRenderSize, fg, bg)
		}

		curr := 0
		o := vr.Begin() + i

		for curr < len(recipe) && (o >= recipe[curr].Region.Begin()) {
			if o < recipe[curr].Region.End() {
				fg = color256(render.Colour(recipe[curr].Flavour.Foreground))
				bg = color256(render.Colour(recipe[curr].Flavour.Background))
			}
			curr++
		}

		iscursor := sel.Contains(Region{o, o})
		if iscursor {
			fg = fg | caretStyle
		}

		if r == '\t' {
			add := (x + 1 + (tabSize - 1)) &^ (tabSize - 1)
			for ; x < add; x++ {
				if x < ex {
					termbox.SetCell(x, y, ' ', fg, bg)
				}
				// A long cursor looks weird
				fg = fg & ^(termbox.AttrUnderline | termbox.AttrReverse)
			}
			continue
		}
		if r == '\n' {
			termbox.SetCell(x, y, ' ', fg, bg)
			x = sx
			y++
			if lineNumbers {
				// This results in additional calls to renderLineNumber.
				// Maybe just accumulate positions needing line numbers, rendering them
				// after the loop?
				renderLineNumber(&line, &x, y, lineNumberRenderSize, defaultFg, defaultBg)
			}
			if y > ey {
				break
			}
			continue
		}
		termbox.SetCell(x, y, r, fg, bg)
		x++
	}
	fg, bg = defaultFg, defaultBg
	// Need this if the cursor is at the end of the buffer
	o := vr.Begin() + len(runes)
	iscursor := sel.Contains(Region{o, o})
	if iscursor {
		fg = fg | caretStyle
		termbox.SetCell(x, y, ' ', fg, bg)
	}

	// restore original caretStyle before blink modification
	caretStyle = oldCaretStyle

	if rs := sel.Regions(); len(rs) > 0 {
		if r := rs[len(rs)-1]; !vr.Covers(r) {
			t.Show(v, r)
		}
	}

	fg, bg = defaultFg, color256(render.Colour{28, 29, 26, 1})
	y = t.window_layout.height - statusbarHeight
	// Draw status bar bottom of window
	for i := 0; i < t.window_layout.width; i++ {
		termbox.SetCell(i, y, ' ', fg, bg)
	}
	go t.renderLStatus(v, y, fg, bg)
	// The right status
	rns := []rune(fmt.Sprintf("Tab Size:%d   %s", tabSize, "Go"))
	x = t.window_layout.width - 1 - len(rns)
	addRunes(x, y, rns, fg, bg)
}

func (t *tbfe) renderLStatus(v *backend.View, y int, fg, bg termbox.Attribute) {
	st := v.Status()
	sel := v.Sel()
	j := 0

	for k, v := range st {
		s := fmt.Sprintf("%s: %s, ", k, v)
		addString(j, y, s, fg, bg)
	}

	if sel.Len() == 0 {
		return
	} else if l := sel.Len(); l > 1 {
		s := fmt.Sprintf("%d selection regions", l)
		j = addString(j, y, s, fg, bg)
	} else if r := sel.Get(0); r.Size() == 0 {
		row, col := v.RowCol(r.A)
		s := fmt.Sprintf("Line %d, Column %d", row+1, col)
		j = addString(j, y, s, fg, bg)
	} else {
		ls := v.Lines(r)
		s := v.Substr(r)
		if len(ls) < 2 {
			s := fmt.Sprintf("%d characters selected", len(s))
			j = addString(j, y, s, fg, bg)
		} else {
			s := fmt.Sprintf("%d lines %d characters selected", len(ls), len(s))
			j = addString(j, y, s, fg, bg)
		}
	}

	if t.status_message != "" {
		s := fmt.Sprintf("; %s", t.status_message)
		addString(j, y, s, fg, bg)
	}
}

func (t *tbfe) clip(v *backend.View, s, e int) Region {
	p := util.Prof.Enter("clip")
	defer p.Exit()
	t.lock.Lock()
	h := t.layout[v].height
	t.lock.Unlock()
	if e-s > h {
		e = s + h
	} else if e-s < h {
		s = e - h
	}
	if e2, _ := v.RowCol(v.TextPoint(e, 0)); e2 < e {
		e = e2
	}
	if s < 0 {
		s = 0
	}
	e = s + h
	r := Region{v.TextPoint(s, 0), v.TextPoint(e, 0)}
	return v.LineR(r)
}

func (t *tbfe) Show(v *backend.View, r Region) {
	t.lock.Lock()
	l := t.layout[v]
	t.lock.Unlock()
	p := util.Prof.Enter("show")
	defer p.Exit()

	lv := l.visible

	s1, _ := v.RowCol(lv.Begin())
	e1, _ := v.RowCol(lv.End())
	s2, _ := v.RowCol(r.Begin())
	e2, _ := v.RowCol(r.End())

	r1 := Region{s1, e1}
	r2 := Region{s2, e2}

	r3 := r1.Cover(r2)
	diff := 0
	if d1, d2 := Abs(r1.Begin()-r3.Begin()), Abs(r1.End()-r3.End()); d1 > d2 {
		diff = r3.Begin() - r1.Begin()
	} else {
		diff = r3.End() - r1.End()
	}
	r3.A = r1.Begin() + diff
	r3.B = r1.End() + diff

	r3 = t.clip(v, r3.A, r3.B)
	l.visible = r3
	t.lock.Lock()
	t.layout[v] = l
	t.lock.Unlock()
	t.render()
}

func (t *tbfe) VisibleRegion(v *backend.View) Region {
	t.lock.Lock()
	r, ok := t.layout[v]
	t.lock.Unlock()
	if !ok || r.lastUpdate != v.ChangeCount() {
		t.Show(v, r.visible)
		t.lock.Lock()
		r = t.layout[v]
		t.lock.Unlock()
	}
	return r.visible
}

func (t *tbfe) StatusMessage(msg string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.status_message = msg
}

func (t *tbfe) ErrorMessage(msg string) {
	log.Error(msg)
}

// TODO(q): Actually show a dialog
func (t *tbfe) MessageDialog(msg string) {
	log.Info(msg)
}

// TODO(q): Actually show a dialog
func (t *tbfe) OkCancelDialog(msg, ok string) bool {
	log.Info(msg, ok)
	return false
}

func (t *tbfe) scroll(b Buffer) {
	t.Show(backend.GetEditor().Console(), Region{b.Size(), b.Size()})
}

func (t *tbfe) Erased(changed_buffer Buffer, region_removed Region, data_removed []rune) {
	t.scroll(changed_buffer)
}

func (t *tbfe) Inserted(changed_buffer Buffer, region_inserted Region, data_inserted []rune) {
	t.scroll(changed_buffer)
}

func (t *tbfe) Prompt(title, folder string, flags int) []string {
	return nil
}

func (t *tbfe) setupCallbacks(view *backend.View) {
	// Ensure that the visible region currently presented is
	// inclusive of the insert/erase delta.
	view.AddObserver(&tbfeBufferDeltaObserver{t: t, view: view})

	backend.OnNew.Add(func(v *backend.View) {
		v.Settings().AddOnChange("lime.frontend.termbox.render", func(name string) { t.render() })
	})

	backend.OnModified.Add(func(v *backend.View) {
		t.render()
	})

	backend.OnSelectionModified.Add(func(v *backend.View) {
		t.render()
	})
}

func (t *tbfe) setupEditor() *backend.Editor {
	ed := backend.GetEditor()

	ed.Init()
	ed.SetDefaultPath("../packages/Default")
	ed.SetUserPath("../packages/User")

	return ed
}

func (t *tbfe) render() {
	t.lock.Lock()
	defer t.lock.Unlock()
	if len(t.dorender) < cap(t.dorender) {
		t.dorender <- true
	}
}

func (t *tbfe) renderthread() {
	pc := 0
	dorender := func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Panic in renderthread: %v\n%s", r, string(debug.Stack()))
				if pc > 1 {
					panic(r)
				}
				pc++
			}
		}()
		termbox.Clear(defaultFg, defaultBg)

		t.lock.Lock()
		vs := make([]*backend.View, 0, len(t.layout))
		ls := make([]layout, 0, len(t.layout))
		for v, l := range t.layout {
			vs = append(vs, v)
			ls = append(ls, l)
		}
		t.lock.Unlock()

		for i, v := range vs {
			t.renderView(v, ls[i])
		}

		termbox.Flush()
	}

	for range t.dorender {
		dorender()
	}
}

func (t *tbfe) handleResize(height, width int, init bool) {
	// This should handle multiple views in a less hardcoded fashion.
	// After all, it is possible to *not* have a view in a window.
	t.lock.Lock()
	if init {
		t.layout[t.currentView] = layout{0, 0, 0, 0, Region{}, 0}
		t.window_layout = layout{0, 0, 0, 0, Region{}, 0}
		t.layout[t.console] = layout{0, 0, 0, 0, Region{}, 0}
	}

	t.window_layout.height = height
	t.window_layout.width = width

	view_layout := t.layout[t.currentView]
	view_layout.height = height - statusbarHeight
	view_layout.width = width
	t.layout[t.currentView] = view_layout
	if *showConsole {
		view_layout := t.layout[t.currentView]
		view_layout.height = height - *consoleHeight - statusbarHeight - 1
		view_layout.width = width

		console_layout := t.layout[t.console]
		console_layout.y = height - *consoleHeight - statusbarHeight
		console_layout.width = width
		console_layout.height = *consoleHeight

		t.layout[t.console] = console_layout
		t.layout[t.currentView] = view_layout
	}
	t.lock.Unlock()

	// Ensure that the new visible region is recalculated
	t.Show(t.currentView, t.VisibleRegion(t.currentView))
}

func (t *tbfe) handleInput(ev termbox.Event) {
	if ev.Key == termbox.KeyCtrlQ {
		t.shutdown <- true
	}

	var kp keys.KeyPress
	if ev.Ch != 0 {
		kp.Key = keys.Key(ev.Ch)
		kp.Text = string(ev.Ch)
	} else if v2, ok := lut[ev.Key]; ok {
		kp = v2
		kp.Text = string(kp.Key)
	} else {
		return
	}

	t.editor.HandleInput(kp)
}

func (t *tbfe) loop() {
	timechan := make(chan bool, 0)

	// Only set up the timers if we should actually blink the cursor
	// This should somehow be changeable on an OnSettingsChanged callback
	if p, _ := t.editor.Settings().Get("caret_blink", true).(bool); p {
		duration := time.Second / 2
		if p, ok := t.editor.Settings().Get("caret_blink_phase", 1.0).(float64); ok {
			duration = time.Duration(float64(time.Second)*p) / 2
		}
		timer := time.NewTimer(duration)

		defer func() {
			timer.Stop()
			close(timechan)
		}()

		go func() {
			for range timer.C {
				timechan <- true
				timer.Reset(duration)
			}
		}()
	}

	// Due to termbox still running, we can't close evchan
	evchan := make(chan termbox.Event, 32)
	go func() {
		for {
			evchan <- termbox.PollEvent()
		}
	}()

	for {
		p := util.Prof.Enter("mainloop")
		select {
		case ev := <-evchan:
			mp := util.Prof.Enter("evchan")
			switch ev.Type {
			case termbox.EventError:
				log.Debug("error occured")
				return
			case termbox.EventResize:
				t.handleResize(ev.Height, ev.Width, false)
			case termbox.EventKey:
				t.handleInput(ev)
				blink = false
			}
			mp.Exit()

		case <-timechan:
			blink = !blink
			t.render()

		case <-t.shutdown:
			return
		}
		p.Exit()
	}
}

func (bdo *tbfeBufferDeltaObserver) Erased(changed_buffer Buffer, region_removed Region, data_removed []rune) {
	ensureVisibleRegionContainsInsertOrEraseDelta(bdo.t, bdo.view, region_removed.A-region_removed.B)
}

func (bdo *tbfeBufferDeltaObserver) Inserted(changed_buffer Buffer, region_inserted Region, data_inserted []rune) {
	ensureVisibleRegionContainsInsertOrEraseDelta(bdo.t, bdo.view, region_inserted.B-region_inserted.A)
}

func ensureVisibleRegionContainsInsertOrEraseDelta(t *tbfe, view *backend.View, delta int) {
	t.lock.Lock()
	visible := t.layout[view].visible
	t.lock.Unlock()
	t.Show(view, Region{visible.Begin(), visible.End() + delta})
}
