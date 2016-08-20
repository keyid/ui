// Copyright (c) 2016 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package ui

import (
	"github.com/richardwilkes/ui/event"
)

// Item represents individual actions that can be issued from a Menu.
type MenuItem struct {
	item          platformMenuItem
	eventHandlers *event.Handlers
	title         string
}

// Title returns this item's title.
func (item *MenuItem) Title() string {
	return item.title
}

// SetKeyModifiers sets the MenuItem's key equivalent modifiers. By default, a MenuItem's modifier
// is set to event.CommandKeyMask.
func (item *MenuItem) SetKeyModifiers(modifierMask event.KeyMask) {
	item.platformSetKeyModifierMask(modifierMask)
}

// SubMenu of this MenuItem or nil.
func (item *MenuItem) SubMenu() *Menu {
	if menu, ok := menuMap[item.platformSubMenu()]; ok {
		return menu
	}
	return nil
}

// EventHandlers implements the event.Target interface.
func (item *MenuItem) EventHandlers() *event.Handlers {
	if item.eventHandlers == nil {
		item.eventHandlers = &event.Handlers{}
	}
	return item.eventHandlers
}

// ParentTarget implements the event.Target interface.
func (item *MenuItem) ParentTarget() event.Target {
	return &App
}
