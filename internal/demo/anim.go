package demo

import "math/rand"

func NewAnimationState(width, height int, seed int64) AnimationState {
	r := rand.New(rand.NewSource(seed))
	count := starCountForArea(width, height)
	stars := make([]Star, 0, count)
	for i := 0; i < count; i++ {
		stars = append(stars, Star{
			X:       r.Intn(width),
			Y:       r.Intn(height),
			Speed:   1 + r.Intn(3),
			Twinkle: r.Intn(10),
		})
	}
	return AnimationState{
		Stars:      stars,
		Meteors:    nil,
		ISSTrail:   nil,
		BootFrames: 14,
	}
}

func ResizeAnimationState(st *AnimationState, width, height int) {
	if width <= 0 || height <= 0 {
		return
	}
	count := starCountForArea(width, height)
	for len(st.Stars) < count {
		i := len(st.Stars)
		st.Stars = append(st.Stars, Star{
			X:       (i*37 + 11) % width,
			Y:       (i*19 + 7) % height,
			Speed:   1 + (i % 3),
			Twinkle: i % 10,
		})
	}
	if len(st.Stars) > count {
		st.Stars = st.Stars[:count]
	}
	for i := range st.Stars {
		if st.Stars[i].X >= width {
			st.Stars[i].X = st.Stars[i].X % width
		}
		if st.Stars[i].Y >= height {
			st.Stars[i].Y = st.Stars[i].Y % height
		}
	}
}

func starCountForArea(width, height int) int {
	area := width * height
	count := area / 24
	if count < 60 {
		count = 60
	}
	if count > 260 {
		count = 260
	}
	return count
}

func AdvanceAnimations(st *AnimationState, width, height, frame int) {
	for i := range st.Stars {
		if frame%st.Stars[i].Speed == 0 {
			st.Stars[i].X--
			if st.Stars[i].X < 0 {
				st.Stars[i].X = width - 1
				st.Stars[i].Y = (st.Stars[i].Y + (i%5 + 1)) % height
			}
		}
		st.Stars[i].Twinkle = (st.Stars[i].Twinkle + 1) % 12
	}

	if frame%9 == 0 {
		for i := range st.Meteors {
			st.Meteors[i].X += st.Meteors[i].DX
			st.Meteors[i].Y += st.Meteors[i].DY
			st.Meteors[i].Life--
		}
	}
	alive := st.Meteors[:0]
	for _, m := range st.Meteors {
		if m.Life > 0 && m.X >= -4 && m.Y >= -2 && m.X < width+4 && m.Y < height+2 {
			alive = append(alive, m)
		}
	}
	st.Meteors = alive

	if frame > 5 && frame%7 == 0 && len(st.Meteors) < 3 {
		spawnMeteor(st, width, frame)
	}

	for i := range st.ISSTrail {
		st.ISSTrail[i].Age++
	}
	trail := st.ISSTrail[:0]
	for _, p := range st.ISSTrail {
		if p.Age < 18 {
			trail = append(trail, p)
		}
	}
	st.ISSTrail = trail
}

func AddISSTrail(st *AnimationState, lat, lon float64) {
	st.ISSTrail = append(st.ISSTrail, TrailPoint{Lat: lat, Lon: lon})
}

func spawnMeteor(st *AnimationState, width, frame int) {
	startX := width - 1 - (frame % 11)
	startY := frame % 6
	st.Meteors = append(st.Meteors, Meteor{
		X:    startX,
		Y:    startY,
		DX:   -2,
		DY:   1,
		Life: 12,
	})
}
