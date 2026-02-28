package demo

import (
	"fmt"
	"math"
	"os"
	"strings"
)

// ColorMode represents terminal color capability.
type ColorMode int

const (
	ColorNone      ColorMode = iota
	ColorTruecolor
)

// CellColor holds an RGB foreground color for a terminal cell.
type CellColor struct {
	R, G, B uint8
	Set     bool
}

// DetectColorMode checks environment variables for truecolor support.
func DetectColorMode(noColor bool) ColorMode {
	if noColor || os.Getenv("NO_COLOR") != "" {
		return ColorNone
	}
	ct := os.Getenv("COLORTERM")
	if ct == "truecolor" || ct == "24bit" {
		return ColorTruecolor
	}
	term := os.Getenv("TERM")
	if strings.Contains(term, "256color") {
		return ColorTruecolor
	}
	tp := os.Getenv("TERM_PROGRAM")
	switch tp {
	case "iTerm.app", "Alacritty", "kitty", "WezTerm":
		return ColorTruecolor
	}
	return ColorNone
}

func makeColorGrid(width, height int) [][]CellColor {
	colors := make([][]CellColor, height)
	for y := range colors {
		colors[y] = make([]CellColor, width)
	}
	return colors
}

func cc(r, g, b uint8) CellColor {
	return CellColor{R: r, G: g, B: b, Set: true}
}

func blendCC(a, b CellColor, t float64) CellColor {
	if !a.Set {
		return b
	}
	if !b.Set {
		return a
	}
	if t <= 0 {
		return a
	}
	if t >= 1 {
		return b
	}
	return CellColor{
		R:   uint8(float64(a.R)*(1-t) + float64(b.R)*t),
		G:   uint8(float64(a.G)*(1-t) + float64(b.G)*t),
		B:   uint8(float64(a.B)*(1-t) + float64(b.B)*t),
		Set: true,
	}
}

func scaleCC(c CellColor, f float64) CellColor {
	if !c.Set {
		return c
	}
	r := float64(c.R) * f
	g := float64(c.G) * f
	b := float64(c.B) * f
	if r > 255 {
		r = 255
	}
	if g > 255 {
		g = 255
	}
	if b > 255 {
		b = 255
	}
	return CellColor{R: uint8(r), G: uint8(g), B: uint8(b), Set: true}
}

// Globe palette
var (
	colOceanDeep  = cc(15, 50, 160)
	colOceanLight = cc(40, 90, 200)
	colLandGreen  = cc(30, 150, 55)
	colLandBrown  = cc(110, 75, 30)
	colNightOcean = cc(5, 10, 40)
	colNightLand  = cc(12, 18, 30)
	colTerminator = cc(210, 110, 40)

	colStarDim    = cc(80, 80, 100)
	colStarMed    = cc(170, 170, 195)
	colStarBright = cc(255, 255, 240)

	colMeteorHead = cc(255, 255, 200)
	colMeteorTail = cc(220, 160, 50)
	colMeteorFade = cc(120, 80, 30)

	colAuroraGreen  = cc(50, 255, 100)
	colAuroraPurple = cc(140, 50, 200)

	colCity      = cc(255, 200, 80)
	colCityPulse = cc(180, 130, 40)

	colISS         = cc(255, 255, 255)
	colISSTrailNew = cc(180, 180, 255)
	colISSTrailOld = cc(70, 70, 110)
)

func globePixelColor(illum, noise, cloud, edgeShade float64) CellColor {
	isLand := noise > 0.07

	var dayCol, nightCol CellColor
	if isLand {
		if noise > 0.15 {
			dayCol = colLandBrown
		} else {
			dayCol = colLandGreen
		}
		nightCol = colNightLand
	} else {
		dayCol = blendCC(colOceanDeep, colOceanLight, (noise+0.5)*0.5)
		nightCol = colNightOcean
	}

	dayFactor := (illum + 1) / 2
	if dayFactor < 0 {
		dayFactor = 0
	}
	if dayFactor > 1 {
		dayFactor = 1
	}
	base := blendCC(nightCol, dayCol, dayFactor)

	if illum > -0.15 && illum < 0.25 {
		tf := 1.0 - math.Abs(illum-0.05)/0.2
		if tf < 0 {
			tf = 0
		}
		if tf > 1 {
			tf = 1
		}
		base = blendCC(base, colTerminator, tf*0.5)
	}

	if cloud > 0.7 && illum > -0.3 {
		cf := (cloud - 0.7) / 0.3
		if cf > 1 {
			cf = 1
		}
		bright := (illum + 1) / 2
		if bright < 0 {
			bright = 0
		}
		cloudCol := cc(uint8(160*bright+40), uint8(170*bright+40), uint8(180*bright+50))
		base = blendCC(base, cloudCol, cf*0.6)
	}

	return scaleCC(base, edgeShade)
}

func buildColoredRow(row []rune, colors []CellColor) string {
	var b strings.Builder
	b.Grow(len(row) * 8)
	var pR, pG, pB uint8
	active := false
	for i, ch := range row {
		if i < len(colors) && colors[i].Set {
			c := colors[i]
			if !active || c.R != pR || c.G != pG || c.B != pB {
				fmt.Fprintf(&b, "\033[38;2;%d;%d;%dm", c.R, c.G, c.B)
				pR, pG, pB = c.R, c.G, c.B
				active = true
			}
		} else if active {
			b.WriteString("\033[0m")
			active = false
		}
		b.WriteRune(ch)
	}
	if active {
		b.WriteString("\033[0m")
	}
	return b.String()
}

func renderColoredFrame(lines []string, width, height int) string {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	blank := strings.Repeat(" ", width)
	var b strings.Builder
	for i := 0; i < height; i++ {
		if i < len(lines) {
			b.WriteString(lines[i])
		} else {
			b.WriteString(blank)
		}
		if i != height-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}
