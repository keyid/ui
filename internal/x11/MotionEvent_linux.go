// Copyright (c) 2016 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package x11

import (
	"github.com/richardwilkes/geom"
	"github.com/richardwilkes/ui/keys"
	// #cgo pkg-config: x11
	// #include <X11/Xlib.h>
	"C"
)

const (
	MotionNotifyType = EventType(C.MotionNotify)
)

type MotionEvent C.XMotionEvent

func (evt *MotionEvent) Window() Window {
	return Window(evt.window)
}

func (evt *MotionEvent) Where() geom.Point {
	return geom.Point{X: float64(evt.x), Y: float64(evt.y)}
}

func (evt *MotionEvent) Modifiers() keys.Modifiers {
	return Modifiers(evt.state)
}
