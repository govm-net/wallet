package screens

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"github.com/lengzhao/wallet/bin/gui/conf"
	"github.com/lengzhao/wallet/bin/gui/res"
)

// customTheme is a simple demonstration of a bespoke theme loaded by a Fyne app.
type customTheme struct {
	background                                                                    color.Color
	button, text, icon, hyperlink, placeholder, primary, hover, scrollBar, shadow color.Color
	regular, bold, italic, bolditalic, monospace                                  fyne.Resource
	disabledButton, disabledIcon, disabledText                                    color.Color
}

func (t *customTheme) BackgroundColor() color.Color {
	return t.background
}

// ButtonColor returns the theme's standard button colour
func (t *customTheme) ButtonColor() color.Color {
	return t.button
}

// DisabledButtonColor returns the theme's disabled button colour
func (t *customTheme) DisabledButtonColor() color.Color {
	return t.disabledButton
}

// HyperlinkColor returns the theme's standard hyperlink colour
func (t *customTheme) HyperlinkColor() color.Color {
	return t.hyperlink
}

// TextColor returns the theme's standard text colour
func (t *customTheme) TextColor() color.Color {
	return t.text
}

// DisabledIconColor returns the color for a disabledIcon UI element
func (t *customTheme) DisabledTextColor() color.Color {
	return t.disabledText
}

// IconColor returns the theme's standard text colour
func (t *customTheme) IconColor() color.Color {
	return t.icon
}

// DisabledIconColor returns the color for a disabledIcon UI element
func (t *customTheme) DisabledIconColor() color.Color {
	return t.disabledIcon
}

// PlaceHolderColor returns the theme's placeholder text colour
func (t *customTheme) PlaceHolderColor() color.Color {
	return t.placeholder
}

// PrimaryColor returns the colour used to highlight primary features
func (t *customTheme) PrimaryColor() color.Color {
	return t.primary
}

// HoverColor returns the colour used to highlight interactive elements currently under a cursor
func (t *customTheme) HoverColor() color.Color {
	return t.hover
}

// FocusColor returns the colour used to highlight a focused widget
func (t *customTheme) FocusColor() color.Color {
	return t.primary
}

// ScrollBarColor returns the color (and translucency) for a scrollBar
func (t *customTheme) ScrollBarColor() color.Color {
	return t.scrollBar
}

// ShadowColor returns the color (and translucency) for shadows used for indicating elevation
func (t *customTheme) ShadowColor() color.Color {
	return t.shadow
}

// TextSize returns the standard text size
func (t *customTheme) TextSize() int {
	return 14
}

func (t *customTheme) TextFont() fyne.Resource {
	if t.regular != nil {
		return t.regular
	}
	return theme.DefaultTextBoldFont()
}

func (t *customTheme) TextBoldFont() fyne.Resource {
	if t.bold != nil {
		return t.bold
	}
	return theme.DefaultTextBoldFont()
}

func (t *customTheme) TextItalicFont() fyne.Resource {
	if t.italic != nil {
		return t.italic
	}
	return theme.DefaultTextBoldItalicFont()
}

func (t *customTheme) TextBoldItalicFont() fyne.Resource {
	if t.bolditalic != nil {
		return t.bolditalic
	}
	return theme.DefaultTextBoldItalicFont()
}

func (t *customTheme) TextMonospaceFont() fyne.Resource {
	if t.regular != nil {
		return t.regular
	}
	return theme.DefaultTextMonospaceFont()
}

func (t *customTheme) Padding() int {
	return 2
}

func (t *customTheme) IconInlineSize() int {
	return 16
}

func (t *customTheme) ScrollBarSize() int {
	return 8
}

func (t *customTheme) ScrollBarSmallSize() int {
	return 5
}

// NewCustomTheme new theme
func NewCustomTheme() fyne.Theme {
	out := &customTheme{
		// background:     color.RGBA{0xf5, 0xf5, 0xf5, 0xff},
		background:     color.RGBA{0xf0, 0xf0, 0xf0, 0xff},
		button:         color.RGBA{0xd9, 0xd9, 0xd9, 0xff},
		disabledButton: color.RGBA{0xe7, 0xe7, 0xe7, 0xff},
		text:           color.RGBA{0x21, 0x21, 0x21, 0xff},
		disabledText:   color.RGBA{0x80, 0x80, 0x80, 0xff},
		icon:           color.RGBA{0x21, 0x21, 0x21, 0xff},
		disabledIcon:   color.RGBA{0x80, 0x80, 0x80, 0xff},
		hyperlink:      color.RGBA{0x0, 0x0, 0xd9, 0xff},
		placeholder:    color.RGBA{0x88, 0x88, 0x88, 0xff},
		primary:        color.RGBA{0x9f, 0xa8, 0xda, 0xff},
		hover:          color.RGBA{0xe7, 0xe7, 0xe7, 0xff},
		scrollBar:      color.RGBA{0x0, 0x0, 0x0, 0x99},
		shadow:         color.RGBA{0x0, 0x0, 0x0, 0x33},
	}

	font := conf.Get().Langure + ".ttf"

	out.regular = res.GetResource(font)
	out.bold = res.GetResource(font)
	out.italic = res.GetResource(font)
	out.bolditalic = res.GetResource(font)

	return out
}
