package event

import (
	"bytes"
	"fmt"
)

// FocusLost is generated when a widget loses the keyboard focus.
type FocusLost struct {
	target   Target
	finished bool
}

// NewFocusLost creates a new FocusLost event. 'target' is the widget that is losing the
// keyboard focus.
func NewFocusLost(target Target) *FocusLost {
	return &FocusLost{target: target}
}

// Type returns the event type ID.
func (e *FocusLost) Type() Type {
	return FocusLostType
}

// Target the original target of the event.
func (e *FocusLost) Target() Target {
	return e.target
}

// Cascade returns true if this event should be passed to its target's parent if not marked done.
func (e *FocusLost) Cascade() bool {
	return false
}

// Finished returns true if this event has been handled and should no longer be processed.
func (e *FocusLost) Finished() bool {
	return e.finished
}

// Finish marks this event as handled and no longer eligible for processing.
func (e *FocusLost) Finish() {
	e.finished = true
}

// String implements the fmt.Stringer interface.
func (e *FocusLost) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("FocusLost[Target: %v", e.target))
	if e.finished {
		buffer.WriteString(", Finished")
	}
	buffer.WriteString("]")
	return buffer.String()
}
