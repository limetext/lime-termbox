package main

import (
	"testing"

	"github.com/limetext/backend/render"
	"github.com/limetext/termbox-go"
)

func TestColourToTermbox(t *testing.T) {
	tests := []struct {
		colour render.Colour
		exp    termbox.Attribute
	}{
		{render.Colour{R: 238, G: 238, B: 238}, termbox.Attribute(255)},
		{render.Colour{R: 255, G: 255, B: 255}, termbox.Attribute(15)},
		{render.Colour{R: 28, G: 29, B: 26}, termbox.Attribute(234)},
	}

	for i, test := range tests {
		if attr := color256(test.colour); attr != test.exp {
			t.Errorf("Test %d: Expected %v got %v", i, test.exp, attr)
		}
	}
}
