package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// GoldenFoxTheme implements the custom theme for Golden Fox Sudoku
type GoldenFoxTheme struct{}

var _ fyne.Theme = (*GoldenFoxTheme)(nil)

// Color scheme
var (
	// Fox orange - primary color
	FoxOrange = color.NRGBA{R: 255, G: 140, B: 0, A: 255}

	// Charcoal - dark backgrounds
	Charcoal = color.NRGBA{R: 54, G: 54, B: 54, A: 255}

	// Soft white
	SoftWhite = color.NRGBA{R: 250, G: 250, B: 250, A: 255}

	// Light gray for cell backgrounds
	LightGray = color.NRGBA{R: 240, G: 240, B: 240, A: 255}

	// Highlight colors
	SelectionHighlight = color.NRGBA{R: 255, G: 200, B: 100, A: 100}
	PeerHighlight      = color.NRGBA{R: 255, G: 220, B: 150, A: 60}
	ConflictRed        = color.NRGBA{R: 255, G: 100, B: 100, A: 255}

	// Given cell color (darker text)
	GivenColor = color.NRGBA{R: 40, G: 40, B: 40, A: 255}

	// User input color (fox orange)
	UserColor = FoxOrange
)

func (g *GoldenFoxTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return FoxOrange
	case theme.ColorNameBackground:
		return SoftWhite
	case theme.ColorNameButton:
		return LightGray
	case theme.ColorNameForeground:
		return Charcoal
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (g *GoldenFoxTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (g *GoldenFoxTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (g *GoldenFoxTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 16
	case theme.SizeNameHeadingText:
		return 24
	default:
		return theme.DefaultTheme().Size(name)
	}
}
