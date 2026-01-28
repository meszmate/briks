package theme

import (
	"fmt"
	"math"

	"github.com/charmbracelet/lipgloss"
)

// RainbowState tracks the current hue offset for the rainbow theme.
type RainbowState struct {
	HueOffset float64
}

// NewRainbowState creates a new rainbow state.
func NewRainbowState() *RainbowState {
	return &RainbowState{}
}

// Tick advances the rainbow hue.
func (r *RainbowState) Tick(delta float64) {
	r.HueOffset += delta * 60 // 60 degrees per second
	if r.HueOffset >= 360 {
		r.HueOffset -= 360
	}
}

// PieceColor returns a cycling color for a piece index (0-6).
func (r *RainbowState) PieceColor(index int) lipgloss.Color {
	hue := math.Mod(r.HueOffset+float64(index)*51.4, 360) // spread 7 pieces across spectrum
	return lipgloss.Color(hslToHex(hue, 0.8, 0.6))
}

// hslToHex converts HSL values to a hex color string.
func hslToHex(h, s, l float64) string {
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	ri := int((r + m) * 255)
	gi := int((g + m) * 255)
	bi := int((b + m) * 255)

	return fmt.Sprintf("#%02x%02x%02x", ri, gi, bi)
}
