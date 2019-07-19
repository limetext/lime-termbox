// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.
package main

import (
	"flag"

	"github.com/limetext/backend/log"
	_ "github.com/limetext/commands"
	py "github.com/limetext/gopy"
	_ "github.com/limetext/sublime"
	"github.com/limetext/util"
	"github.com/nsf/termbox-go"
)

// Command line flags
var (
	showConsole   = flag.Bool("console", false, "Display console")
	consoleHeight = flag.Int("consoleHeight", 20, "Height of console")
	rotateLog     = flag.Bool("rotateLog", false, "Rotate debug log")
)

func main() {
	flag.Parse()

	log.AddFilter("file", log.FINEST, log.NewFileLogWriter("debug.log", *rotateLog))
	// Replace Global Logger filter so that it does not interfere with the ui
	log.AddFilter("stdout", log.DEBUG, log.NewFileLogWriter("debug.log", *rotateLog))

	if err := termbox.Init(); err != nil {
		log.Error(err)
		return
	}

	defer shutdown()

	t := createFrontend()
	go t.renderthread()
	t.loop()
}

func shutdown() {
	defer log.Close()

	py.NewLock()
	py.Finalize()

	termbox.Close()

	log.Debug(util.Prof)
	if err := recover(); err != nil {
		log.Critical(err)
		panic(err)
	}
}
