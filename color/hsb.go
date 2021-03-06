package color

import (
	"math"
)

// HSB creates a new opaque Color from HSB values in the range 0-1.
func HSB(hue, saturation, brightness float64) Color {
	return HSBA(hue, saturation, brightness, 1)
}

// HSBA creates a new Color from HSBA values in the range 0-1.
func HSBA(hue, saturation, brightness, alpha float64) Color {
	saturation = clamp0To1(saturation)
	brightness = clamp0To1(brightness)
	v := clamp0To1AndScale255(brightness)
	if saturation == 0 {
		return RGBA(v, v, v, alpha)
	}
	h := (hue - math.Floor(hue)) * 6
	f := h - math.Floor(h)
	p := clamp0To1AndScale255(brightness * (1 - saturation))
	q := clamp0To1AndScale255(brightness * (1 - saturation*f))
	t := clamp0To1AndScale255(brightness * (1 - (saturation * (1 - f))))
	switch int(h) {
	case 0:
		return RGBA(v, t, p, alpha)
	case 1:
		return RGBA(q, v, p, alpha)
	case 2:
		return RGBA(p, v, t, alpha)
	case 3:
		return RGBA(p, q, v, alpha)
	case 4:
		return RGBA(t, p, v, alpha)
	default:
		return RGBA(v, p, q, alpha)
	}
}

// Hue of the color, a value from 0-1.
func (c Color) Hue() float64 {
	hue, _, _ := c.HSB()
	return hue
}

// SetHue creates a new color from this color with the specified hue, a value from 0-1.
func (c Color) SetHue(hue float64) Color {
	_, s, b := c.HSB()
	return HSBA(hue, s, b, c.AlphaIntensity())
}

// AdjustHue creates a new color from this color with its hue adjusted by the specified amount.
func (c Color) AdjustHue(amount float64) Color {
	h, s, b := c.HSB()
	return HSBA(h+amount, s, b, c.AlphaIntensity())
}

// Saturation of the color, a value from 0-1.
func (c Color) Saturation() float64 {
	brightness := c.Brightness()
	if brightness != 0 {
		return (brightness - (float64(min(c.Red(), c.Green(), c.Blue())) / 255)) / brightness
	}
	return 0
}

// SetSaturation creates a new color from this color with the specified saturation.
func (c Color) SetSaturation(saturation float64) Color {
	h, _, b := c.HSB()
	return HSBA(h, saturation, b, c.AlphaIntensity())
}

// AdjustSaturation creates a new color from this color with its saturation adjusted by the
// specified amount.
func (c Color) AdjustSaturation(amount float64) Color {
	h, s, b := c.HSB()
	return HSBA(h, s+amount, b, c.AlphaIntensity())
}

// Brightness of the color, a value from 0-1.
func (c Color) Brightness() float64 {
	return float64(max(c.Red(), c.Green(), c.Blue())) / 255
}

// SetBrightness creates a new color from this color with the specified brightness.
func (c Color) SetBrightness(brightness float64) Color {
	h, s, _ := c.HSB()
	return HSBA(h, s, brightness, c.AlphaIntensity())
}

// AdjustBrightness creates a new color from this color with its brightness adjusted by the
// specified amount.
func (c Color) AdjustBrightness(amount float64) Color {
	h, s, b := c.HSB()
	return HSBA(h, s, b+amount, c.AlphaIntensity())
}

// HSB returns the hue, saturation and brightness of the color. Values are in the range 0-1.
func (c Color) HSB() (hue, saturation, brightness float64) {
	r := c.Red()
	g := c.Green()
	b := c.Blue()
	cmax := max(r, g, b)
	cmin := min(r, g, b)
	brightness = float64(cmax) / 255
	if cmax != 0 {
		saturation = float64(cmax-cmin) / float64(cmax)
	} else {
		saturation = 0
	}
	if saturation == 0 {
		hue = 0
	} else {
		div := float64(cmax - cmin)
		rc := float64(cmax-r) / div
		gc := float64(cmax-g) / div
		bc := float64(cmax-b) / div
		if r == cmax {
			hue = bc - gc
		} else if g == cmax {
			hue = 2 + rc - bc
		} else {
			hue = 4 + gc - rc
		}
		hue /= 6
		if hue < 0 {
			hue++
		}
	}
	return
}

func clamp0To1(value float64) float64 {
	switch {
	case value < 0:
		return 0
	case value > 1:
		return 1
	default:
		return value
	}
}

func clamp0To255(value int) int {
	switch {
	case value < 0:
		return 0
	case value > 255:
		return 255
	default:
		return value
	}
}

func clamp0To1AndScale255(value float64) int {
	return clamp0To255(int(clamp0To1(value)*255 + 0.5))
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func max(a, b, c int) int {
	if a > b {
		if a > c {
			return a
		}
		return c
	}
	if b > c {
		return b
	}
	return c
}
