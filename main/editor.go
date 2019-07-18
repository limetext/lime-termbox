package main

import (
	"github.com/limetext/backend"
	"github.com/limetext/backend/log"
)

func setSchemeSettings(ed *backend.Editor) {
	s := ed.GetColorScheme(ed.Settings().Get("color_scheme", "").(string))
	if s == nil {
		log.Error("No colour scheme to set defaults from")
		return
	}

	defaultFg = color256(s.GlobalSettings().Foreground)
	defaultBg = color256(s.GlobalSettings().Background)
}

func createNewView(filename string, window *backend.Window) *backend.View {
	v := window.OpenFile(filename, 0)

	v.Settings().Set("trace", true)

	return v
}
