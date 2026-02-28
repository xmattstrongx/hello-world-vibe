package demo

import (
	"math"
	"strings"
)

type planet struct {
	Name     string
	Char     rune
	Orbit    float64 // relative orbital radius (0-1)
	Speed    float64 // orbital speed multiplier
	HasRing  bool
	Label    bool // whether to label this planet
}

var planets = []planet{
	{Name: "Mercury", Char: '.', Orbit: 0.12, Speed: 4.15},
	{Name: "Venus", Char: 'o', Orbit: 0.20, Speed: 1.62},
	{Name: "Earth", Char: 'O', Orbit: 0.28, Speed: 1.00, Label: true},
	{Name: "Mars", Char: '*', Orbit: 0.36, Speed: 0.53},
	{Name: "Jupiter", Char: '#', Orbit: 0.52, Speed: 0.08},
	{Name: "Saturn", Char: '@', Orbit: 0.68, Speed: 0.03, HasRing: true},
	{Name: "Uranus", Char: '~', Orbit: 0.82, Speed: 0.01},
	{Name: "Neptune", Char: ':', Orbit: 0.94, Speed: 0.006},
}

const zoomTransitionFrames = 12

func advanceZoom(zoom *ZoomState, wantSolar bool) {
	if wantSolar {
		zoom.Active = true
		zoom.Progress += 1.0 / float64(zoomTransitionFrames)
		if zoom.Progress > 1.0 {
			zoom.Progress = 1.0
		}
	} else {
		zoom.Progress -= 1.0 / float64(zoomTransitionFrames)
		if zoom.Progress <= 0 {
			zoom.Progress = 0
			zoom.Active = false
		}
	}
}

func renderSolarSystem(grid [][]rune, frame int, zoom ZoomState) {
	height := len(grid)
	if height == 0 {
		return
	}
	width := len(grid[0])

	cx := float64(width) / 2
	cy := float64(height) / 2

	// Aspect ratio compensation for monospace
	aspect := 2.2

	// Max orbital radius in characters
	maxRX := float64(width)/2 - 2
	maxRY := float64(height)/2 - 1

	// Draw orbit paths (dotted ellipses)
	for _, p := range planets {
		orbitRX := p.Orbit * maxRX
		orbitRY := p.Orbit * maxRY
		steps := int(math.Max(float64(width), float64(height)) * p.Orbit * 3)
		if steps < 40 {
			steps = 40
		}
		for i := 0; i < steps; i++ {
			angle := 2 * math.Pi * float64(i) / float64(steps)
			x := int(cx + orbitRX*math.Cos(angle))
			y := int(cy + orbitRY*math.Sin(angle)/aspect)
			if x >= 0 && x < width && y >= 0 && y < height {
				if grid[y][x] == ' ' && i%4 == 0 {
					grid[y][x] = '·'
				}
			}
		}
	}

	// Draw Sun at center
	sunChars := []rune{'*', 'O', '*'}
	for i, ch := range sunChars {
		sx := int(cx) - 1 + i
		sy := int(cy)
		if sx >= 0 && sx < width && sy >= 0 && sy < height {
			grid[sy][sx] = ch
		}
	}
	// Sun glow
	glowOffsets := [][2]int{{0, -1}, {0, 1}, {-2, 0}, {2, 0}}
	for _, off := range glowOffsets {
		gx := int(cx) + off[0]
		gy := int(cy) + off[1]
		if gx >= 0 && gx < width && gy >= 0 && gy < height {
			if grid[gy][gx] == ' ' || grid[gy][gx] == '·' {
				grid[gy][gx] = '+'
			}
		}
	}

	// Draw planets at their orbital positions
	t := float64(frame) * 0.02
	for _, p := range planets {
		orbitRX := p.Orbit * maxRX
		orbitRY := p.Orbit * maxRY
		angle := t * p.Speed
		px := int(cx + orbitRX*math.Cos(angle))
		py := int(cy + orbitRY*math.Sin(angle)/aspect)

		if px >= 0 && px < width && py >= 0 && py < height {
			grid[py][px] = p.Char

			// Draw Saturn's ring
			if p.HasRing {
				ringOffsets := [][2]int{{-2, 0}, {-1, 0}, {1, 0}, {2, 0}}
				for _, ro := range ringOffsets {
					rx := px + ro[0]
					ry := py + ro[1]
					if rx >= 0 && rx < width && ry >= 0 && ry < height && grid[ry][rx] != p.Char {
						grid[ry][rx] = '-'
					}
				}
			}

			// Label Earth
			if p.Label {
				label := "Earth"
				lx := px - len(label)/2
				ly := py - 1
				if ly < 0 {
					ly = py + 1
				}
				if ly >= 0 && ly < height {
					for i, ch := range label {
						xx := lx + i
						if xx >= 0 && xx < width {
							grid[ly][xx] = ch
						}
					}
				}
			}
		}
	}
}

func renderTransitionFrame(data FrameData, termW, termH int) string {
	progress := data.Zoom.Progress

	if progress >= 1.0 {
		return renderSolarFrame(data, termW, termH)
	}
	if progress <= 0 {
		return RenderFrame(data, termW, termH)
	}

	// During transition: render solar system with a shrinking earth indicator
	return renderSolarFrame(data, termW, termH)
}

func renderSolarFrame(data FrameData, termW, termH int) string {
	view := data.ViewMode
	if view == "" {
		view = resolveViewMode(termW, termH, "auto")
	}
	_, _, leftPad := globeViewport(termW, termH, view)

	// Use most of the terminal for the solar system
	sWidth := termW - 4
	if sWidth < 40 {
		sWidth = 40
	}
	sHeight := termH - 6
	if sHeight < 12 {
		sHeight = 12
	}

	grid := makeGrid(sWidth, sHeight)

	// Draw background stars
	if data.Anim != nil {
		// Reuse starfield but adapted to solar system grid
		for _, s := range data.Anim.Stars {
			sx := s.X % sWidth
			sy := s.Y % sHeight
			if sx >= 0 && sx < sWidth && sy >= 0 && sy < sHeight && grid[sy][sx] == ' ' {
				phase := (s.Twinkle + data.Frame) % 12
				ch := '.'
				if phase > 8 {
					ch = '+'
				}
				if phase == 0 {
					ch = '*'
				}
				grid[sy][sx] = ch
			}
		}
	}

	renderSolarSystem(grid, data.Frame, data.Zoom)

	// Build output
	_ = leftPad
	solarPad := (termW - sWidth) / 2
	if solarPad < 0 {
		solarPad = 0
	}
	pad := strings.Repeat(" ", solarPad)

	lines := make([]string, 0, termH)
	quality := "std"
	if data.MaxASCII {
		quality = "max"
	}
	header := fitLine("SOLAR SYSTEM VIEW | UTC "+data.Now.Format("15:04:05")+" | ASCII "+quality+" | Press 'z' to return", termW)
	lines = append(lines, header)
	border := strings.Repeat("=", sWidth)
	lines = append(lines, pad+border)
	for _, row := range grid {
		lines = append(lines, pad+string(row))
	}
	lines = append(lines, pad+border)

	// Info line
	lines = append(lines, fitLine("Planets orbit the Sun · Earth highlighted · Press 'z' to zoom back to Earth", termW))
	lines = append(lines, fitLine("Controls: z zoom | [space]/p pause | q quit", termW))

	return renderFixedFrame(lines, termW, termH)
}
