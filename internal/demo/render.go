package demo

import "math"

func makeGrid(width, height int) [][]rune {
	grid := make([][]rune, height)
	for y := range grid {
		grid[y] = make([]rune, width)
		for x := range grid[y] {
			grid[y][x] = ' '
		}
	}
	return grid
}

func renderGlobe(grid [][]rune, subLon, decl float64, frame int, view string, maxASCII bool) [][]bool {
	height := len(grid)
	if height == 0 {
		return nil
	}
	width := len(grid[0])
	mask := make([][]bool, height)
	for y := range mask {
		mask[y] = make([]bool, width)
	}
	cx := float64(width-1) / 2
	cy := float64(height-1) / 2
	rx := float64(width-2) / 2
	ry := float64(height-2) / 2
	ramp := []rune(" .,:-~=+*#%@")
	if view == "compact" {
		ramp = []rune(" .,:-+*#@")
	}
	if !maxASCII {
		ramp = []rune(" ;:,.+*#@")
	}
	if maxASCII && view != "compact" {
		ramp = []rune(" .'`^,:;~-_=+*#%@")
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx := (float64(x) - cx) / rx
			dy := (float64(y) - cy) / ry
			r2 := dx*dx + dy*dy
			if r2 > 1 {
				continue
			}
			mask[y][x] = true

			z := math.Sqrt(1 - r2)
			lat := math.Asin(dy) * 180 / math.Pi
			lon := math.Atan2(dx, z)*180/math.Pi + subLon
			illum := SolarIllumination(lat, lon, decl, subLon)
			if illum < -1 {
				illum = -1
			}
			if illum > 1 {
				illum = 1
			}
			edgeShade := z*0.55 + 0.45
			noise := terrainNoise(lat, lon, float64(frame))
			cloud := cloudNoise(lat, lon, float64(frame))
			terrainBias := 0.0
			if maxASCII {
				if noise > 0.07 {
					terrainBias = 0.05
				} else {
					terrainBias = -0.03
				}
				terrainBias += (cloud - 0.5) * 0.05
			}
			level := (illum*0.5 + 0.5 + terrainBias) * edgeShade
			if level < 0 {
				level = 0
			}
			if level > 1 {
				level = 1
			}
			idx := int(level * float64(len(ramp)-1))
			if idx < 0 {
				idx = 0
			}
			if idx >= len(ramp) {
				idx = len(ramp) - 1
			}
			grid[y][x] = ramp[idx]
			if maxASCII && noise > 0.18 && illum > -0.2 && idx > 2 {
				grid[y][x] = '#'
			}
		}
	}
	return mask
}

func terrainNoise(lat, lon, frame float64) float64 {
	l := lat * math.Pi / 180
	o := lon * math.Pi / 180
	return 0.42*math.Sin(o*2.3+l*1.7) + 0.33*math.Cos(o*4.9-l*2.2) + 0.25*math.Sin((o+l)*3.7+frame*0.03)
}

func cloudNoise(lat, lon, frame float64) float64 {
	l := lat * math.Pi / 180
	o := lon * math.Pi / 180
	v := math.Sin(o*7.2+frame*0.12)*math.Cos(l*6.1-frame*0.04) + math.Sin((o-l)*5.3)
	return (v + 2) / 4
}

func drawStarfield(grid [][]rune, stars []Star, frame int, mask [][]bool) {
	for _, s := range stars {
		if s.Y < 0 || s.Y >= len(grid) || s.X < 0 || s.X >= len(grid[0]) {
			continue
		}
		if mask[s.Y][s.X] {
			continue
		}
		phase := (s.Twinkle + frame) % 12
		ch := '.'
		if phase > 8 {
			ch = '+'
		}
		if phase == 0 {
			ch = '*'
		}
		grid[s.Y][s.X] = ch
	}
}

func drawMeteors(grid [][]rune, meteors []Meteor, mask [][]bool) {
	for _, m := range meteors {
		x := m.X
		y := m.Y
		for tail := 0; tail < 4; tail++ {
			px := x + tail
			py := y - tail/2
			if py < 0 || py >= len(grid) || px < 0 || px >= len(grid[0]) {
				continue
			}
			if mask[py][px] {
				continue
			}
			ch := '\\'
			if tail > 1 {
				ch = '-'
			}
			if tail == 3 {
				ch = '.'
			}
			grid[py][px] = ch
		}
	}
}

func drawAurora(grid [][]rune, mask [][]bool, subLon, decl float64, frame int) {
	width := len(grid[0])
	height := len(grid)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if !mask[y][x] {
				continue
			}
			lat, lon, _, ok := pixelToWorld(x, y, width, height, subLon)
			if !ok || math.Abs(lat) < 52 {
				continue
			}
			illum := SolarIllumination(lat, lon, decl, subLon)
			if illum > 0.15 {
				continue
			}
			wave := math.Sin((float64(x)+float64(frame)*0.8)/3.6) + math.Cos((float64(y)+float64(frame)*0.3)/2.1)
			if wave > 1.0 {
				grid[y][x] = '^'
			} else if wave > 0.6 {
				grid[y][x] = '~'
			}
		}
	}
}

func plotCities(grid [][]rune, cities []City, subLon float64) {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return
	}
	width := len(grid[0])
	height := len(grid)
	for _, c := range cities {
		x, y, ok := Project(c.Lat, c.Lon, subLon, width, height)
		if ok {
			grid[y][x] = 'o'
		}
	}
}

func drawCityPulses(grid [][]rune, cities []City, subLon float64, frame int) {
	width := len(grid[0])
	height := len(grid)
	for i, c := range cities {
		x, y, ok := Project(c.Lat, c.Lon, subLon, width, height)
		if !ok {
			continue
		}
		phase := (frame + i*3) % 12
		if phase > 4 {
			continue
		}
		r := phase
		if r == 0 {
			continue
		}
		pts := [][2]int{
			{x - r, y}, {x + r, y}, {x, y - r}, {x, y + r},
		}
		for _, p := range pts {
			if p[1] >= 0 && p[1] < height && p[0] >= 0 && p[0] < width && grid[p[1]][p[0]] == ' ' {
				grid[p[1]][p[0]] = '.'
			}
		}
	}
}

func drawISSTrail(grid [][]rune, trail []TrailPoint, subLon float64) {
	width := len(grid[0])
	height := len(grid)
	for _, p := range trail {
		x, y, ok := Project(p.Lat, p.Lon, subLon, width, height)
		if !ok {
			continue
		}
		ch := '.'
		if p.Age < 12 {
			ch = ':'
		}
		if p.Age < 5 {
			ch = '*'
		}
		if grid[y][x] != 'o' {
			grid[y][x] = ch
		}
	}
}

func plotISS(grid [][]rune, iss ISSMarker, subLon float64) {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return
	}
	x, y, ok := Project(iss.Lat, iss.Lon, subLon, len(grid[0]), len(grid))
	if ok {
		grid[y][x] = '@'
	}
}

func applyScanlines(grid [][]rune, frame int) {
	for y := 0; y < len(grid); y++ {
		if (y+frame)%4 != 0 {
			continue
		}
		for x := 0; x < len(grid[0]); x++ {
			if grid[y][x] == '*' || grid[y][x] == '@' || grid[y][x] == 'o' {
				continue
			}
			if grid[y][x] == '#' {
				grid[y][x] = '*'
			} else if grid[y][x] == '+' {
				grid[y][x] = '.'
			}
		}
	}
}

func pixelToWorld(x, y, width, height int, subLon float64) (lat, lon, z float64, ok bool) {
	cx := float64(width-1) / 2
	cy := float64(height-1) / 2
	rx := float64(width-2) / 2
	ry := float64(height-2) / 2
	dx := (float64(x) - cx) / rx
	dy := (float64(y) - cy) / ry
	r2 := dx*dx + dy*dy
	if r2 > 1 {
		return 0, 0, 0, false
	}
	z = math.Sqrt(1 - r2)
	lat = math.Asin(dy) * 180 / math.Pi
	lon = math.Atan2(dx, z)*180/math.Pi + subLon
	return lat, lon, z, true
}
