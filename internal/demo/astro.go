package demo

import (
	"math"
	"sort"
	"time"
)

func SolarDeclination(t time.Time) float64 {
	day := float64(t.YearDay())
	return 23.44 * math.Sin((2*math.Pi/365.0)*(day-81))
}

func SubsolarLongitude(t time.Time) float64 {
	h := float64(t.Hour()) + float64(t.Minute())/60 + float64(t.Second())/3600
	lon := -15 * (h - 12)
	for lon < -180 {
		lon += 360
	}
	for lon > 180 {
		lon -= 360
	}
	return lon
}

func IsDaylight(lat, lon, decl, subLon float64) bool {
	return SolarIllumination(lat, lon, decl, subLon) > 0
}

func SolarIllumination(lat, lon, decl, subLon float64) float64 {
	latR := degToRad(lat)
	lonR := degToRad(lon)
	declR := degToRad(decl)
	subLonR := degToRad(subLon)
	cosZenith := math.Sin(latR)*math.Sin(declR) + math.Cos(latR)*math.Cos(declR)*math.Cos(lonR-subLonR)
	return cosZenith
}

func DaylightCities(lang string, decl, subLon float64) []City {
	list := make([]City, 0, len(Cities))
	for _, c := range Cities {
		if IsDaylight(c.Lat, c.Lon, decl, subLon) {
			x := c
			if lang == "en" {
				x.Greeting = "Hello"
			}
			list = append(list, x)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Lon < list[j].Lon
	})
	return list
}

func Project(lat, lon, centerLon float64, width, height int) (int, int, bool) {
	latR := degToRad(lat)
	lonR := degToRad(lon - centerLon)
	x := math.Cos(latR) * math.Sin(lonR)
	y := math.Sin(latR)
	z := math.Cos(latR) * math.Cos(lonR)
	if z <= 0 {
		return 0, 0, false
	}

	cx := float64(width-1) / 2
	cy := float64(height-1) / 2
	rx := float64(width-2) / 2
	ry := float64(height-2) / 2
	px := int(cx + x*rx)
	py := int(cy + y*ry)
	if px < 0 || px >= width || py < 0 || py >= height {
		return 0, 0, false
	}
	return px, py, true
}

func degToRad(v float64) float64 {
	return v * math.Pi / 180
}
