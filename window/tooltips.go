package window

import (
	"time"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ui"
	"github.com/richardwilkes/ui/layout"
)

var (
	// TooltipDelay holds the delay before a tooltip will be shown.
	TooltipDelay = 1500 * time.Millisecond
	// TooltipDismissal holds the delay before a tooltip will be dismissed.
	TooltipDismissal = 3 * time.Second
)

type tooltipSequencer struct {
	window   *Window
	avoid    geom.Rect
	sequence int
}

func (ts *tooltipSequencer) show() {
	if ts.window.tooltipSequence == ts.sequence {
		tip := ts.window.lastToolTip
		_, pref, _ := ui.Sizes(tip, layout.NoHintSize)
		bounds := geom.Rect{Point: geom.Point{X: ts.avoid.X, Y: ts.avoid.Y + ts.avoid.Height + 1}, Size: pref}
		if bounds.X < 0 {
			bounds.X = 0
		}
		if bounds.Y < 0 {
			bounds.Y = 0
		}
		viewSize := ts.window.root.Size()
		if viewSize.Width < bounds.Width {
			_, pref, _ := ui.Sizes(tip, geom.Size{Width: viewSize.Width, Height: layout.NoHint})
			if viewSize.Width < pref.Width {
				bounds.X = 0
				bounds.Width = viewSize.Width
			} else {
				bounds.Width = pref.Width
			}
			bounds.Height = pref.Height
		}
		if viewSize.Width < bounds.X+bounds.Width {
			bounds.X = viewSize.Width - bounds.Width
		}
		if viewSize.Height < bounds.Y+bounds.Height {
			bounds.Y = ts.avoid.Y - (bounds.Height + 1)
			if bounds.Y < 0 {
				bounds.Y = 0
			}
		}
		tip.SetBounds(bounds)
		ts.window.root.SetTooltip(tip)
		ts.window.lastTooltipShownAt = time.Now()
		ts.window.InvokeAfter(ts.close, TooltipDismissal)
	}
}

func (ts *tooltipSequencer) close() {
	if ts.window.tooltipSequence == ts.sequence {
		ts.window.root.SetTooltip(nil)
	}
}
