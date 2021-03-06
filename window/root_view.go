package window

import (
	"fmt"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ui"
	"github.com/richardwilkes/ui/color"
	"github.com/richardwilkes/ui/layout"
	"github.com/richardwilkes/ui/menu"
	"github.com/richardwilkes/ui/widget"
)

// RootView provides a root view for a window.
type RootView struct {
	widget.Block
	tooltip ui.Widget
	menuBar menu.Bar
	content ui.Widget
}

func newRootView(window ui.Window) *RootView {
	view := &RootView{}
	view.SetBackground(color.Background)
	view.SetWindow(window)
	view.Describer = func() string { return fmt.Sprintf("RootView #%d", view.ID()) }
	view.SetLayout(&RootLayout{view: view})
	view.content = widget.NewBlock()
	view.AddChild(view.content)
	return view
}

// MenuBar returns the menu bar.
func (view *RootView) MenuBar() menu.Bar {
	return view.menuBar
}

// SetMenuBar sets the menu bar.
func (view *RootView) SetMenuBar(bar menu.Bar) {
	if view.menuBar != nil {
		if actual, ok := view.menuBar.(ui.Widget); ok {
			view.RemoveChild(actual)
		}
	}
	view.menuBar = bar
	if actual, ok := bar.(ui.Widget); ok {
		view.AddChildAtIndex(actual, 0)
	}
}

// Tooltip returns the tooltip for this component.
func (view *RootView) Tooltip() ui.Widget {
	return view.tooltip
}

// SetTooltip sets the tooltip for this component.
func (view *RootView) SetTooltip(tip ui.Widget) {
	if view.tooltip != nil {
		view.tooltip.Repaint()
		view.RemoveChild(view.tooltip)
	}
	view.tooltip = tip
	if tip != nil {
		view.AddChild(tip)
		tip.Repaint()
	}
}

// Content returns the content area.
func (view *RootView) Content() ui.Widget {
	return view.content
}

// RootLayout holds layout information.
type RootLayout struct {
	view *RootView
}

// Sizes implements the Sizer interface.
func (lay *RootLayout) Sizes(hint geom.Size) (min, pref, max geom.Size) {
	min, pref, max = ui.Sizes(lay.view.content, hint)
	if lay.view.menuBar != nil {
		if actual, ok := lay.view.menuBar.(ui.Widget); ok {
			_, barSize, _ := ui.Sizes(actual, layout.NoHintSize)
			lay.adjustSizeForBarSize(&min, barSize)
			lay.adjustSizeForBarSize(&pref, barSize)
			lay.adjustSizeForBarSize(&max, barSize)
		}
	}
	return
}

func (lay *RootLayout) adjustSizeForBarSize(size *geom.Size, barSize geom.Size) {
	size.Height += barSize.Height
	if size.Width < barSize.Width {
		size.Width = barSize.Width
	}
}

// Layout implements the Layout interface.
func (lay *RootLayout) Layout() {
	bounds := lay.view.LocalBounds()
	if lay.view.menuBar != nil {
		if actual, ok := lay.view.menuBar.(ui.Widget); ok {
			_, size, _ := ui.Sizes(actual, layout.NoHintSize)
			actual.SetBounds(geom.Rect{Size: geom.Size{Width: bounds.Width, Height: size.Height}})
			bounds.Y += size.Height
			bounds.Height -= size.Height
		}
	}
	lay.view.content.SetBounds(bounds)
}
