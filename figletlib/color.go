package figletlib

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// RGB represents a color with red, green, blue components (0-255).
type RGB struct {
	R, G, B uint8
}

// ColorMode represents different coloring modes.
type ColorMode int

const (
	ColorModeNone ColorMode = iota
	ColorModeGradient
	ColorModeRainbow
)

// ColorConfig holds color configuration.
type ColorConfig struct {
	Mode       ColorMode
	StartColor RGB
	EndColor   RGB
}

// ANSI escape codes for colors.
const (
	Reset = "\033[0m"
)

// ParseColor parses a color string in various formats:
// - hex: #FF0000, #ff0000, FF0000, ff0000
// - rgb: rgb(255,0,0), RGB(255,0,0)
// - named colors: red, green, blue, etc.
func ParseColor(colorStr string) (RGB, error) {
	colorStr = strings.TrimSpace(colorStr)

	// Handle hex colors
	colorStr = strings.TrimPrefix(colorStr, "#")

	if len(colorStr) == 6 {
		// Parse hex color
		r, err := strconv.ParseUint(colorStr[0:2], 16, 8)
		if err != nil {
			return RGB{}, fmt.Errorf("invalid hex color: %s", colorStr)
		}
		g, err := strconv.ParseUint(colorStr[2:4], 16, 8)
		if err != nil {
			return RGB{}, fmt.Errorf("invalid hex color: %s", colorStr)
		}
		b, err := strconv.ParseUint(colorStr[4:6], 16, 8)
		if err != nil {
			return RGB{}, fmt.Errorf("invalid hex color: %s", colorStr)
		}

		return RGB{uint8(r), uint8(g), uint8(b)}, nil
	}

	// Handle rgb() format
	if strings.HasPrefix(strings.ToLower(colorStr), "rgb(") && strings.HasSuffix(colorStr, ")") {
		rgbStr := colorStr[4 : len(colorStr)-1]
		parts := strings.Split(rgbStr, ",")
		if len(parts) != 3 {
			return RGB{}, fmt.Errorf("invalid rgb format: %s", colorStr)
		}

		r, err := strconv.ParseUint(strings.TrimSpace(parts[0]), 10, 8)
		if err != nil {
			return RGB{}, fmt.Errorf("invalid rgb red value: %s", parts[0])
		}
		g, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 8)
		if err != nil {
			return RGB{}, fmt.Errorf("invalid rgb green value: %s", parts[1])
		}
		b, err := strconv.ParseUint(strings.TrimSpace(parts[2]), 10, 8)
		if err != nil {
			return RGB{}, fmt.Errorf("invalid rgb blue value: %s", parts[2])
		}

		return RGB{uint8(r), uint8(g), uint8(b)}, nil
	}

	// Handle named colors
	namedColors := map[string]RGB{
		"red":     {255, 0, 0},
		"green":   {0, 255, 0},
		"blue":    {0, 0, 255},
		"yellow":  {255, 255, 0},
		"magenta": {255, 0, 255},
		"cyan":    {0, 255, 255},
		"white":   {255, 255, 255},
		"black":   {0, 0, 0},
		"orange":  {255, 165, 0},
		"purple":  {128, 0, 128},
		"pink":    {255, 192, 203},
		"lime":    {0, 255, 0},
		"navy":    {0, 0, 128},
		"teal":    {0, 128, 128},
		"silver":  {192, 192, 192},
		"gray":    {128, 128, 128},
		"maroon":  {128, 0, 0},
		"olive":   {128, 128, 0},
	}

	if color, exists := namedColors[strings.ToLower(colorStr)]; exists {
		return color, nil
	}

	return RGB{}, fmt.Errorf("unknown color: %s", colorStr)
}

// ToANSI converts RGB to ANSI escape sequence for 24-bit color.
func (rgb RGB) ToANSI() string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", rgb.R, rgb.G, rgb.B)
}

// Interpolate creates a color between two colors based on factor (0.0 to 1.0).
func (start RGB) Interpolate(end RGB, factor float64) RGB {
	if factor < 0 {
		factor = 0
	}
	if factor > 1 {
		factor = 1
	}

	r := float64(start.R) + (float64(end.R)-float64(start.R))*factor
	g := float64(start.G) + (float64(end.G)-float64(start.G))*factor
	b := float64(start.B) + (float64(end.B)-float64(start.B))*factor

	return RGB{uint8(r), uint8(g), uint8(b)}
}

// HSVtoRGB converts HSV to RGB.
func HSVtoRGB(h, s, v float64) RGB {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	var r, g, b float64

	switch {
	case h >= 0 && h < 60:
		r, g, b = c, x, 0
	case h >= 60 && h < 120:
		r, g, b = x, c, 0
	case h >= 120 && h < 180:
		r, g, b = 0, c, x
	case h >= 180 && h < 240:
		r, g, b = 0, x, c
	case h >= 240 && h < 300:
		r, g, b = x, 0, c
	case h >= 300 && h < 360:
		r, g, b = c, 0, x
	}

	return RGB{
		uint8((r + m) * 255),
		uint8((g + m) * 255),
		uint8((b + m) * 255),
	}
}

// GetRainbowColor returns a rainbow color based on position (0.0 to 1.0).
func GetRainbowColor(position float64) RGB {
	// Ensure position is in range [0, 1]
	if position < 0 {
		position = 0
	}
	if position > 1 {
		position = 1
	}

	// Convert position to hue (0-360 degrees)
	hue := position * 360

	// Use full saturation and value for vibrant colors
	return HSVtoRGB(hue, 1.0, 1.0)
}

// ApplyColor applies color to a character based on its position and configuration.
func ApplyColor(char rune, position float64, totalWidth int, config ColorConfig) string {
	if config.Mode == ColorModeNone {
		return string(char)
	}

	var color RGB

	switch config.Mode {
	case ColorModeGradient:
		color = config.StartColor.Interpolate(config.EndColor, position)
	case ColorModeRainbow:
		color = GetRainbowColor(position)
	case ColorModeNone:
		return string(char)
	default:
		return string(char)
	}

	return color.ToANSI() + string(char) + Reset
}
