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
	"github.com/richardwilkes/ui/draw"
)

// ImageLabel represents a non-interactive image.
type ImageLabel struct {
	Block
	image *draw.Image
}

// NewImageLabel creates an ImageLabel with the specified image.
func NewImageLabel(img *draw.Image) *ImageLabel {
	return NewImageLabelWithImageSize(img, draw.Size{})
}

// NewImageLabelWithImageSize creates a new ImageLabel with the specified image. The image will be
// set to the specified size.
func NewImageLabelWithImageSize(img *draw.Image, size draw.Size) *ImageLabel {
	label := &ImageLabel{image: img}
	if size.Width <= 0 || size.Height <= 0 {
		label.SetSizer(label)
	} else {
		label.SetSizer(&imageLabelSizer{label: label, size: size})
	}
	label.AddEventHandler(PaintEvent, label.paint)
	return label
}

// Sizes implements Sizer
func (label *ImageLabel) Sizes(hint draw.Size) (min, pref, max draw.Size) {
	size := label.image.Size()
	if border := label.Border(); border != nil {
		size.AddInsets(border.Insets())
	}
	return size, size, size
}

func (label *ImageLabel) paint(event *Event) {
	bounds := label.LocalInsetBounds()
	size := label.image.Size()
	if size.Width < bounds.Width {
		bounds.X += (bounds.Width - size.Width) / 2
		bounds.Width = size.Width
	}
	if size.Height < bounds.Height {
		bounds.Y += (bounds.Height - size.Height) / 2
		bounds.Height = size.Height
	}
	event.GC.DrawImageInRect(label.image, bounds)
}

type imageLabelSizer struct {
	label *ImageLabel
	size  draw.Size
}

// Sizes implements Sizer
func (sizer *imageLabelSizer) Sizes(hint draw.Size) (min, pref, max draw.Size) {
	pref = sizer.size
	if border := sizer.label.Border(); border != nil {
		pref.AddInsets(border.Insets())
	}
	return pref, pref, pref
}
