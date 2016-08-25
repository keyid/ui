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
	"github.com/richardwilkes/ui/event"
)

// Item represents individual actions that can be issued from a Menu.
type Item struct {
	item          PlatformItem
	eventHandlers *event.Handlers
	title         string
}

// PlatformPtr returns the underlying platform data pointer.
func (item *Item) PlatformPtr() PlatformItem {
	return item.item
}

// Title returns this item's title.
func (item *Item) Title() string {
	return item.title
}

// SetKeyModifiers sets the menu item's key equivalent modifiers. By default, a menu item's modifier
// is set to event.CommandKeyMask.
func (item *Item) SetKeyModifiers(modifierMask event.KeyMask) {
	item.platformSetKeyModifierMask(modifierMask)
}

// SubMenu of this menu item or nil.
func (item *Item) SubMenu() *Menu {
	if menu, ok := menuMap[item.platformSubMenu()]; ok {
		return menu
	}
	return nil
}

// EventHandlers implements the event.Target interface.
func (item *Item) EventHandlers() *event.Handlers {
	if item.eventHandlers == nil {
		item.eventHandlers = &event.Handlers{}
	}
	return item.eventHandlers
}

// ParentTarget implements the event.Target interface.
func (item *Item) ParentTarget() event.Target {
	return ParentTarget()
}