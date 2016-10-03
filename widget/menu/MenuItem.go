// Copyright (c) 2016 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package menu

import (
	"fmt"
	"github.com/richardwilkes/geom"
	"github.com/richardwilkes/ui/color"
	"github.com/richardwilkes/ui/event"
	"github.com/richardwilkes/ui/keys"
	"github.com/richardwilkes/ui/widget"
	"unicode/utf8"
)

type MenuItem struct {
	widget.Block
	Theme        *Theme
	Title        string
	KeyCode      int
	KeyModifiers keys.Modifiers
	menu         *Menu
	pos          float64
	highlighted  bool
	menuOpen     bool
}

func NewMenuItem(title string, keyCode int, handler event.Handler) *MenuItem {
	return NewMenuItemWithModifiers(title, keyCode, keys.PlatformMenuModifier(), handler)
}

func NewMenuItemWithModifiers(title string, keyCode int, modifiers keys.Modifiers, handler event.Handler) *MenuItem {
	item := &MenuItem{}
	item.Theme = StdTheme
	item.Title = title
	item.KeyCode = keyCode
	item.KeyModifiers = modifiers
	item.Describer = func() string { return fmt.Sprintf("MenuItem #%d (%s)", item.ID(), item.Title) }
	item.SetSizer(item)
	handlers := item.EventHandlers()
	handlers.Add(event.PaintType, item.paint)
	handlers.Add(event.MouseDownType, item.mouseDown)
	handlers.Add(event.MouseDraggedType, item.mouseDragged)
	handlers.Add(event.MouseUpType, item.mouseUp)
	handlers.Add(event.MouseEnteredType, item.mouseEntered)
	handlers.Add(event.MouseMovedType, item.mouseMoved)
	handlers.Add(event.MouseExitedType, item.mouseExited)
	handlers.Add(event.KeyDownType, item.keyDown)
	if handler != nil {
		handlers.Add(event.SelectionType, handler)
	}
	return item
}

// Sizes implements Sizer
func (item *MenuItem) Sizes(hint geom.Size) (min, pref, max geom.Size) {
	pref = item.Theme.TitleFont.Measure(item.Title)
	pref.Width += item.Theme.HMargin*2 + item.Theme.KeySpacing
	pref.Height += item.Theme.VMargin * 2
	pref.GrowToInteger()
	if item.KeyCode != 0 {
		mapping := keys.MappingForKeyCode(item.KeyCode)
		if mapping != nil {
			keySize := item.Theme.KeyFont.Measure(mapping.Name)
			pref.Width += keySize.Width
			if pref.Height < keySize.Height {
				pref.Height = keySize.Height
			}
			mods := item.KeyModifiers
			if mods == 0 {
				mods = keys.PlatformMenuModifier()
			}
			modSize := item.Theme.KeyFont.Measure(mods.String())
			modSize.GrowToInteger()
			pref.Width += modSize.Width
		}
		pref.GrowToInteger()
	}
	pref.ConstrainForHint(hint)
	if border := item.Border(); border != nil {
		pref.AddInsets(border.Insets())
	}
	return pref, pref, pref
}

func (item *MenuItem) calculateAcceleratorPosition() float64 {
	pos := item.Theme.HMargin
	if item.KeyCode != 0 {
		mapping := keys.MappingForKeyCode(item.KeyCode)
		if mapping != nil {
			pos += item.Theme.KeyFont.Measure(mapping.Name).Width
		}
	}
	return pos
}

func (item *MenuItem) paint(evt event.Event) {
	bounds := item.LocalInsetBounds()
	gc := evt.(*event.Paint).GC()
	gc.SetColor(item.currentBackground())
	gc.FillRect(bounds)
	size := item.Theme.TitleFont.Measure(item.Title)
	gc.SetColor(item.textColor())
	gc.DrawString(bounds.X+item.Theme.HMargin, bounds.Y+(bounds.Height-size.Height)/2, item.Title, item.Theme.TitleFont)
	if item.KeyCode != 0 && item.pos > 0 {
		mapping := keys.MappingForKeyCode(item.KeyCode)
		if mapping != nil {
			size = item.Theme.KeyFont.Measure(mapping.Name)
			x := bounds.X + bounds.Width - item.pos
			y := bounds.Y + (bounds.Height-size.Height)/2
			modY := y - (item.Theme.KeyFont.Leading() + 0.5)
			needOffset := false
			for _, r := range mapping.Name {
				if r < ' ' || r > utf8.RuneSelf {
					needOffset = true
				} else {
					needOffset = false
					break
				}
			}
			if needOffset {
				y = modY
			}
			gc.DrawString(x, y, mapping.Name, item.Theme.KeyFont)
			modText := item.KeyModifiers.String()
			size = item.Theme.KeyFont.Measure(modText)
			gc.DrawString(x-size.Width, modY, modText, item.Theme.KeyFont)
		}
	}
}

func (item *MenuItem) currentBackground() color.Color {
	switch {
	case !item.Enabled():
		return item.Theme.Background.AdjustBrightness(item.Theme.DisabledAdjustment)
	case item.highlighted || item.menuOpen:
		return item.Theme.HighlightedBackground
	case item.Focused():
		return item.Theme.Background.Blend(color.KeyboardFocus, 0.5)
	default:
		return item.Theme.Background
	}
}

func (item *MenuItem) textColor() color.Color {
	if !item.Enabled() {
		return item.Theme.TextWhenDisabled
	}
	if item.currentBackground().Luminance() > 0.65 {
		return item.Theme.TextWhenLight
	}
	return item.Theme.TextWhenDark
}

func (item *MenuItem) mouseDown(evt event.Event) {
	item.highlighted = true
	item.Repaint()
}

func (item *MenuItem) mouseOver(where geom.Point) {
	if item.Enabled() {
		bounds := item.LocalInsetBounds()
		highlighted := bounds.Contains(item.FromWindow(where))
		if item.highlighted != highlighted {
			item.highlighted = highlighted
			item.Repaint()
		}
	}
}

func (item *MenuItem) mouseDragged(evt event.Event) {
	item.mouseOver(evt.(*event.MouseDragged).Where())
}

func (item *MenuItem) mouseUp(evt event.Event) {
	item.highlighted = false
	item.Repaint()
	bounds := item.LocalInsetBounds()
	mouseUp := evt.(*event.MouseUp)
	if bounds.Contains(item.FromWindow(mouseUp.Where())) {
		event.Dispatch(event.NewClosing(item))
		event.Dispatch(event.NewSelection(item))
	}
}

func (item *MenuItem) mouseEntered(evt event.Event) {
	item.mouseOver(evt.(*event.MouseEntered).Where())
}

func (item *MenuItem) mouseMoved(evt event.Event) {
	item.mouseOver(evt.(*event.MouseMoved).Where())
}

func (item *MenuItem) mouseExited(evt event.Event) {
	item.mouseOver(evt.(*event.MouseExited).Where())
}

func (item *MenuItem) keyDown(evt event.Event) {
	if keys.IsControlAction(evt.(*event.KeyDown).Code()) {
		evt.Finish()
		event.Dispatch(event.NewSelection(item))
	}
}

// SubMenu of this menu item or nil.
func (item *MenuItem) SubMenu() *Menu {
	return item.menu
}