// Copyright (c) 2016 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package widget

import (
	"github.com/richardwilkes/geom"
	"github.com/richardwilkes/ui/cursor"
	"github.com/richardwilkes/ui/internal/iwindow"
	"time"
	"unsafe"
	// #cgo linux LDFLAGS: -lX11 -lcairo
	// #include <stdlib.h>
	// #include <stdio.h>
	// #include <string.h>
	// #include <X11/Xlib.h>
	// #include <cairo/cairo.h>
	// #include <cairo/cairo-xlib.h>
	// #include "Types.h"
	"C"
)

var (
	lastKnownWindowBounds = make(map[platformWindow]geom.Rect)
)

func toXWindow(window platformWindow) C.Window {
	return C.Window(uintptr(window))
}

func toXDrawable(window platformWindow) C.Drawable {
	return C.Drawable(uintptr(window))
}

func toPlatformWindow(window C.Window) platformWindow {
	return platformWindow(uintptr(window))
}

func platformGetKeyWindow() platformWindow {
	var focus C.Window
	var revert C.int
	C.XGetInputFocus(xDisplay, &focus, &revert)
	return toPlatformWindow(focus)
}

func platformBringAllWindowsToFront() {
	// RAW: Implement for Linux
}

func platformHideCursorUntilMouseMoves() {
	// RAW: Implement for Linux
}

func platformNewWindow(bounds geom.Rect, styleMask WindowStyleMask) (window platformWindow, surface platformSurface) {
	screen := C.XDefaultScreen(xDisplay)
	var windowAttributes C.XSetWindowAttributes
	windowAttributes.background_pixmap = C.None
	windowAttributes.backing_store = C.WhenMapped
	win := C.XCreateWindow(xDisplay, C.XRootWindow(xDisplay, screen), C.int(bounds.X), C.int(bounds.Y), C.uint(bounds.Width), C.uint(bounds.Height), 0, C.CopyFromParent, C.InputOutput, nil, C.CWBackPixmap|C.CWBackingStore, &windowAttributes)
	lastKnownWindowBounds[toPlatformWindow(win)] = bounds
	C.XSelectInput(xDisplay, win, C.KeyPressMask|C.KeyReleaseMask|C.ButtonPressMask|C.ButtonReleaseMask|C.EnterWindowMask|C.LeaveWindowMask|C.ExposureMask|C.PointerMotionMask|C.ExposureMask|C.VisibilityChangeMask|C.StructureNotifyMask|C.FocusChangeMask)
	C.XSetWMProtocols(xDisplay, win, &wmDeleteAtom, C.True)
	iwindow.Count++
	return toPlatformWindow(win), platformSurface(C.cairo_xlib_surface_create(xDisplay, C.Drawable(uintptr(win)), C.XDefaultVisual(xDisplay, screen), C.int(bounds.Width), C.int(bounds.Height)))
}

func (window *Wnd) platformClose() {
	delete(lastKnownWindowBounds, window.window)
	iwindow.Count--
	C.cairo_surface_destroy(window.surface)
	C.XDestroyWindow(xDisplay, toXWindow(window.window))
}

func (window *Wnd) platformTitle() string {
	var result *C.char
	C.XFetchName(xDisplay, toXWindow(window.window), &result)
	if result == nil {
		return ""
	}
	defer C.XFree(result)
	return C.GoString(result)
}

func (window *Wnd) platformSetTitle(title string) {
	cTitle := C.CString(title)
	C.XStoreName(xDisplay, toXWindow(window.window), cTitle)
	C.free(unsafe.Pointer(cTitle))
}

func (window *Wnd) platformFrame() geom.Rect {
	// Use the last set bounds instead of querying the server. I do this because reporting often lags behind, which
	// means the call to XGetGeometry may not have the correct values.
	if bounds, ok := lastKnownWindowBounds[window.window]; ok {
		return bounds
	}
	var root C.Window
	var x, y C.int
	var width, height, border, depth C.uint
	C.XGetGeometry(xDisplay, toXDrawable(window.window), &root, &x, &y, &width, &height, &border, &depth)
	return geom.Rect{Point: geom.Point{X: float64(x), Y: float64(y)}, Size: geom.Size{Width: float64(width), Height: float64(height)}}
}

func (window *Wnd) platformSetFrame(bounds geom.Rect) {
	lastKnownWindowBounds[window.window] = bounds
	win := toXWindow(window.window)
	C.XMoveResizeWindow(xDisplay, win, C.int(bounds.X), C.int(bounds.Y), C.uint(bounds.Width), C.uint(bounds.Height))
}

func (window *Wnd) platformContentFrame() geom.Rect {
	// RAW: Implement for Linux
	return window.platformFrame()
}

func (window *Wnd) platformToFront() {
	win := toXWindow(window.window)
	if window.wasMapped {
		C.XRaiseWindow(xDisplay, win)
	} else {
		window.wasMapped = true
		bounds, ok := lastKnownWindowBounds[window.window]
		C.XMapWindow(xDisplay, win)
		if ok {
			C.XMoveWindow(xDisplay, win, C.int(bounds.X), C.int(bounds.Y))
		}
	}
}

func (window *Wnd) platformRepaint(bounds geom.Rect) {
	event := C.XExposeEvent{_type: C.Expose, window: toXWindow(window.window), x: C.int(bounds.X), y: C.int(bounds.Y), width: C.int(bounds.Width), height: C.int(bounds.Height)}
	C.XSendEvent(xDisplay, toXWindow(window.window), 0, C.ExposureMask, (*C.XEvent)(unsafe.Pointer(&event)))
}

func (window *Wnd) platformFlushPainting() {
	C.XFlush(xDisplay)
}

func (window *Wnd) platformScalingFactor() float64 {
	// RAW: Implement for Linux
	return 1
}

func (window *Wnd) platformMinimize() {
	C.XIconifyWindow(xDisplay, toXWindow(window.window), C.XDefaultScreen(xDisplay))
}

func (window *Wnd) platformZoom() {
	// RAW: Implement for Linux
}

func (window *Wnd) platformSetToolTip(tip string) {
	// RAW: Implement for Linux
}

func (window *Wnd) platformSetCursor(c *cursor.Cursor) {
	// RAW: Implement for Linux
}

func (window *Wnd) platformInvoke(id uint64) {
	if window.Valid() {
		event := C.XClientMessageEvent{_type: C.ClientMessage, message_type: goTaskAtom, format: 32}
		data := (*uint64)(unsafe.Pointer(&event.data))
		*data = id
		C.XSendEvent(xDisplay, toXWindow(window.window), 0, C.NoEventMask, (*C.XEvent)(unsafe.Pointer(&event)))
		C.XFlush(xDisplay)
	}
}

func (window *Wnd) platformInvokeAfter(id uint64, after time.Duration) {
	time.AfterFunc(after, func() {
		window.platformInvoke(id)
	})
}