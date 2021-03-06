package event

import (
	"bytes"
	"fmt"

	"github.com/richardwilkes/ui/keys"
)

// KeyDown is generated when a key is pressed.
type KeyDown struct {
	target    Target
	code      int
	modifiers keys.Modifiers
	ch        rune
	repeat    bool
	finished  bool
	discarded bool
}

// NewKeyDown creates a new KeyDown event. 'target' is the widget that has the keyboard focus.
// 'code' is the virtual key code. 'ch' is the rune (may be 0). 'modifiers' are the keyboard
// modifiers keys that were down. 'autoRepeat' is true if the key is auto-repeating.
func NewKeyDown(target Target, code int, ch rune, modifiers keys.Modifiers, autoRepeat bool) *KeyDown {
	return &KeyDown{target: target, code: code, ch: ch, modifiers: modifiers, repeat: autoRepeat}
}

// Type returns the event type ID.
func (e *KeyDown) Type() Type {
	return KeyDownType
}

// Target the original target of the event.
func (e *KeyDown) Target() Target {
	return e.target
}

// Cascade returns true if this event should be passed to its target's parent if not marked done.
func (e *KeyDown) Cascade() bool {
	return true
}

// Finished returns true if this event has been handled and should no longer be processed.
func (e *KeyDown) Finished() bool {
	return e.finished
}

// Finish marks this event as handled and no longer eligible for processing.
func (e *KeyDown) Finish() {
	e.finished = true
}

// Code returns the virtual key code.
func (e *KeyDown) Code() int {
	return e.code
}

// Rune returns the rune that was typed. May be 0.
func (e *KeyDown) Rune() rune {
	return e.ch
}

// Modifiers returns the key modifiers that were down.
func (e *KeyDown) Modifiers() keys.Modifiers {
	return e.modifiers
}

// Repeat returns true if this key was generated as part of an auto-repeating key.
func (e *KeyDown) Repeat() bool {
	return e.repeat
}

// Discarded returns true if this event should be treated as if it never happened.
func (e *KeyDown) Discarded() bool {
	return e.discarded
}

// Discard marks this event to be thrown away.
func (e *KeyDown) Discard() {
	e.discarded = true
	e.finished = true
}

// String implements the fmt.Stringer interface.
func (e *KeyDown) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("KeyDown[")
	if e.discarded {
		buffer.WriteString("Discarded, ")
	}
	buffer.WriteString(fmt.Sprintf("Code: %d", e.code))
	if e.ch != 0 {
		buffer.WriteString(fmt.Sprintf(", Rune: %d (%s)", e.ch, string(e.ch)))
	}
	buffer.WriteString(fmt.Sprintf(", Target: %v", e.target))
	modifiers := e.modifiers.String()
	if modifiers != "" {
		buffer.WriteString(", ")
		buffer.WriteString(modifiers)
	}
	if e.repeat {
		buffer.WriteString(", Auto-Repeat")
	}
	if e.finished {
		buffer.WriteString(", Finished")
	}
	buffer.WriteString("]")
	return buffer.String()
}
