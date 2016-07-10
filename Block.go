// Copyright (c) 2016 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package ui

// Block is the basic graphical block in a window.
type Block struct {
	Layout
	Border
	// Called when no layout has been set for this block. Returns the minimum, preferred, and
	// maximum sizes of the block. The hint's values will be either NoLayoutHint or a specific value
	// if that particular dimension has already been determined.
	Sizes func(hint Size) (min, pref, max Size)
	// Called to handle mouse down events on this block. Return true if the mouse down event was
	// passed off to a popup menu.
	OnMouseDown func(where Point, keyModifiers int, which int, clickCount int) bool
	// Called to handle mouse dragged events on this block.
	OnMouseDragged func(where Point, keyModifiers int)
	// Called to handle mouse up events on this block.
	OnMouseUp func(where Point, keyModifiers int)
	// Called to handle mouse entered events on this block.
	OnMouseEntered func(where Point, keyModifiers int)
	// Called to handle mouse moved events on this block.
	OnMouseMoved func(where Point, keyModifiers int)
	// Called to handle mouse exited events on this block.
	OnMouseExited func(where Point, keyModifiers int)
	// Called when a tooltip is being requested for the block at the specified position.
	OnToolTip func(where Point) string
	// Called to draw the block's contents.
	OnPaint    func(g Graphics, dirty Rect)
	id         int
	bounds     Rect
	window     *Window
	parent     *Block
	children   []*Block
	layoutData interface{}
	background Color
	NeedLayout bool
	focused    bool
	disabled   bool
}

var (
	nextBlockID = 1
	// DefaultMinSize is the minimum size value that will be used for blocks that don't have a layout and don't provide a Sizes function.
	DefaultMinSize = Size{Width: 0, Height: 0}
	// DefaultPrefSize is the preferred size value that will be used for blocks that don't have a layout and don't provide a Sizes function.
	DefaultPrefSize = Size{Width: 0, Height: 0}
)

// NewBlock creates a new, empty block.
func NewBlock() *Block {
	b := &Block{}
	b.Init()
	return b
}

// Init initializes the block.
func (b *Block) Init() {
	b.id = nextBlockID
	nextBlockID++
}

// ID returns the id for this block.
func (b *Block) ID() int {
	return b.id
}

// Children returns the blocks contained by this block.
func (b *Block) Children() []*Block {
	return b.children
}

// LayoutData returns any layout data that is associated with this block.
func (b *Block) LayoutData() interface{} {
	return b.layoutData
}

// SetLayoutData sets the layout data for this block.
func (b *Block) SetLayoutData(data interface{}) {
	if b.layoutData != data {
		b.layoutData = data
		b.NeedLayout = true
	}
}

// IndexOfChild returns the index of the specified child, or -1 if the passed in block is not a
// child of this block.
func (b *Block) IndexOfChild(child *Block) int {
	for i, one := range b.children {
		if one == child {
			return i
		}
	}
	return -1
}

// AddChild adds the specified block as a child of this block, removing it from any previous
// parent it may have had.
func (b *Block) AddChild(child *Block) {
	child.RemoveFromParent()
	child.parent = b
	b.children = append(b.children, child)
	b.NeedLayout = true
}

// RemoveFromParent removes this block from its parent, if any.
func (b *Block) RemoveFromParent() {
	if b.parent != nil {
		i := b.parent.IndexOfChild(b)
		copy(b.parent.children[i:], b.parent.children[i+1:])
		length := len(b.parent.children) - 1
		b.parent.children[length] = nil
		b.parent.children = b.parent.children[:length]
		b.parent.NeedLayout = true
		b.parent = nil
	}
}

// Location returns the location of this block in its parent's coordinate system.
func (b *Block) Location() Point {
	return b.bounds.Point
}

// SetLocation sets the location of this block in its parent's coordinate system.
func (b *Block) SetLocation(pt Point) {
	if b.bounds.Point != pt {
		b.Repaint()
		b.bounds.Point = pt
		b.Repaint()
	}
}

// Size returns the size of this block.
func (b *Block) Size() Size {
	return b.bounds.Size
}

// SetSize sets the size of this block.
func (b *Block) SetSize(size Size) {
	if b.bounds.Size != size {
		b.Repaint()
		b.bounds.Size = size
		b.NeedLayout = true
		b.Repaint()
	}
}

// Insets returns the margins of this block, as determined by any border that it may have.
func (b *Block) Insets() Insets {
	if b.Border != nil {
		return b.Border.Insets()
	}
	return Insets{}
}

// Bounds returns the location and size of this block in its parent's coordinate system.
func (b *Block) Bounds() Rect {
	return b.bounds
}

// LocalBounds returns the location and size of this block in local coordinates.
func (b *Block) LocalBounds() Rect {
	return b.bounds.CopyAndZeroLocation()
}

// LocalInsetBounds returns the location and size of this block in local coordinates after
// applying its insets.
func (b *Block) LocalInsetBounds() Rect {
	bounds := b.LocalBounds()
	bounds.Inset(b.Insets())
	return bounds
}

// SetBounds sets the location and size of this block in its parent's coordinate system.
func (b *Block) SetBounds(bounds Rect) {
	moved := b.bounds.X != bounds.X || b.bounds.Y != bounds.Y
	resized := b.bounds.Width != bounds.Width || b.bounds.Height != bounds.Height
	if moved || resized {
		b.Repaint()
		if moved {
			b.bounds.Point = bounds.Point
		}
		if resized {
			b.bounds.Size = bounds.Size
			b.NeedLayout = true
		}
		b.Repaint()
	}
}

// ComputeSizes is called by the Layout interface to determine the size constraints for
// this block. This method will defer to a layout that has been set on the block. If no layout
// has been set, then it will attempt to obtain the sizes by calling Sizes(hint). If no Sizes
// function has been set, then a default set of sizes are returned.
func (b *Block) ComputeSizes(hint Size) (min, pref, max Size) {
	if b.Layout != nil {
		return b.Layout.ComputeSizes(b, hint)
	}
	if b.Sizes != nil {
		return b.Sizes(hint)
	}
	return DefaultMinSize, DefaultPrefSize, DefaultLayoutMaxSize(DefaultPrefSize)
}

func (b *Block) paint(g Graphics, dirty Rect) {
	dirty.Intersect(b.LocalBounds())
	if !dirty.IsEmpty() {
		b.paintSelf(g, dirty)
		for _, child := range b.children {
			adjusted := dirty
			adjusted.Intersect(child.bounds)
			if !adjusted.IsEmpty() {
				b.paintChild(child, g, adjusted)
			}
		}
	}
}

func (b *Block) paintSelf(g Graphics, dirty Rect) {
	g.Save()
	defer g.Restore()
	g.ClipRect(dirty)
	if b.background.Alpha() > 0 {
		g.SetFillColor(b.background)
		g.FillRect(dirty)
	}
	b.paintBorder(g)
	if b.OnPaint != nil {
		b.OnPaint(g, dirty)
	}
}

func (b *Block) paintBorder(g Graphics) {
	if b.Border != nil {
		g.Save()
		defer g.Restore()
		b.Border.Paint(g, b.LocalBounds())
	}
}

func (b *Block) paintChild(child *Block, g Graphics, dirty Rect) {
	g.Save()
	defer g.Restore()
	g.Translate(child.bounds.X, child.bounds.Y)
	dirty.X -= child.bounds.X
	dirty.Y -= child.bounds.Y
	child.paint(g, dirty)
}

// Background returns the current background color of this block.
func (b *Block) Background() Color {
	return b.background
}

// SetBackground sets the background color of this block.
func (b *Block) SetBackground(color Color) {
	if color != b.background {
		b.background = color
		b.Repaint()
	}
}

// Repaint marks this block for painting at the next update.
func (b *Block) Repaint() {
	b.RepaintBounds(b.LocalBounds())
}

// RepaintBounds marks the specified bounds within the block for painting at the next update.
func (b *Block) RepaintBounds(bounds Rect) {
	bounds.Intersect(b.LocalBounds())
	if !bounds.IsEmpty() {
		if b.parent != nil {
			bounds.X += b.bounds.X
			bounds.Y += b.bounds.Y
			b.parent.RepaintBounds(bounds)
		} else if b.window != nil {
			b.window.RepaintBounds(bounds)
		}
	}
}

// ValidateLayout triggers any layout that needs to be run by this block or its children.
func (b *Block) ValidateLayout() {
	if b.NeedLayout {
		if b.Layout != nil {
			b.Layout.Layout(b)
			b.Repaint()
		}
		b.NeedLayout = false
	}
	for _, child := range b.children {
		child.ValidateLayout()
	}
}

// BlockAt returns the leaf-most child block containing the location, or this block if no child
// is found.
func (b *Block) BlockAt(pt Point) *Block {
	for _, child := range b.children {
		if child.bounds.Contains(pt) {
			pt.Subtract(child.bounds.Point)
			return child.BlockAt(pt)
		}
	}
	return b
}

// ToWindow converts block-local coordinates into window coordinates.
func (b *Block) ToWindow(pt Point) Point {
	pt.Add(b.bounds.Point)
	parent := b.parent
	for parent != nil {
		pt.Add(parent.bounds.Point)
		parent = parent.parent
	}
	return pt
}

// FromWindow converts window coordinates into block-local coordinates.
func (b *Block) FromWindow(pt Point) Point {
	pt.Subtract(b.bounds.Point)
	parent := b.parent
	for parent != nil {
		pt.Subtract(parent.bounds.Point)
		parent = parent.parent
	}
	return pt
}

// Window returns the containing window, or nil.
func (b *Block) Window() *Window {
	if b.window != nil {
		return b.window
	}
	if b.parent != nil {
		return b.parent.Window()
	}
	return nil
}

// InLiveResize returns true if the window is being actively resized by the user at this point in
// time. If it is, expensive painting operations should be deferred if possible to give a smooth
// resizing experience.
func (b *Block) InLiveResize() bool {
	w := b.Window()
	if w != nil {
		return w.inLiveResize
	}
	return false
}

// Opaque returns true if this block's background is fully opaque.
func (b *Block) Opaque() bool {
	return b.background.Opaque()
}

// Enabled returns true if this block is currently enabled.
func (b *Block) Enabled() bool {
	return !b.disabled
}

// Disabled returns true if this block is currently disabled.
func (b *Block) Disabled() bool {
	return b.disabled
}

// SetDisabled sets this block's enabled state.
func (b *Block) SetDisabled(disabled bool) {
	if b.disabled != disabled {
		b.disabled = disabled
		b.Repaint()
	}
}

// Focused returns true if this block has the keyboard focus.
func (b *Block) Focused() bool {
	return b.focused && !b.disabled
}

// SetFocused sets this block's focus state.
func (b *Block) SetFocused(focused bool) {
	if b.focused != focused {
		b.focused = focused
		b.Repaint()
	}
}

// ToolTip returns the tooltip for the specified point within this block by calling the OnToolTip()
// function if it is set, otherwise it returns an empty string.
func (b *Block) ToolTip(where Point) string {
	if b.OnToolTip != nil {
		return b.OnToolTip(where)
	}
	return ""
}