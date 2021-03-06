package textfield

import (
	"time"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ui/border"
	"github.com/richardwilkes/ui/color"
	"github.com/richardwilkes/ui/font"
)

var (
	// StdTheme is the theme all new TextFields get by default.
	StdTheme = NewTheme()
)

// Theme contains the theme elements for TextFields.
type Theme struct {
	Font                    *font.Font    // The font to use.
	Border                  border.Border // The border to use when not focused.
	FocusBorder             border.Border // The border to use when focused.
	BlinkRate               time.Duration // The rate at which the cursor blinks.
	MinimumTextWidth        float64       // The minimum space to permit for text.
	DisabledBackgroundColor color.Color   // The color to use for the background when disabled.
	InvalidBackgroundColor  color.Color   // The color to use for the background when marked invalid.
}

// NewTheme creates a new TextField theme.
func NewTheme() *Theme {
	theme := &Theme{}
	theme.Init()
	return theme
}

// Init initializes the theme with its default values.
func (theme *Theme) Init() {
	theme.Font = font.User
	theme.Border = border.NewCompound(border.NewLine(color.Background.AdjustBrightness(-0.25), geom.NewUniformInsets(1)), border.NewEmpty(geom.Insets{Top: 1, Left: 4, Bottom: 1, Right: 4}))
	theme.FocusBorder = border.NewCompound(border.NewLine(color.KeyboardFocus, geom.NewUniformInsets(2)), border.NewEmpty(geom.Insets{Top: 0, Left: 3, Bottom: 0, Right: 3}))
	theme.BlinkRate = time.Millisecond * 560
	theme.MinimumTextWidth = 10
	theme.DisabledBackgroundColor = color.Background
	theme.InvalidBackgroundColor = color.RGB(255, 232, 232)
}
