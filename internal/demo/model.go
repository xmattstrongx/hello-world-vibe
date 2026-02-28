package demo

import "time"

type City struct {
	Name     string
	Country  string
	Lat      float64
	Lon      float64
	Greeting string
}

type IPGeoResponse struct {
	City string  `json:"city"`
	Lat  float64 `json:"latitude"`
	Lon  float64 `json:"longitude"`
}

type SunriseResponse struct {
	Results struct {
		Sunrise string `json:"sunrise"`
		Sunset  string `json:"sunset"`
	} `json:"results"`
	Status string `json:"status"`
}

type ISSMarker struct {
	Lat float64
	Lon float64
	OK  bool
}

type FrameData struct {
	Now         time.Time
	Lang        string
	ViewMode    string
	MaxASCII    bool
	Frame       int
	RawInput    bool
	Controls    ControlState
	Zoom        ZoomState
	Local       *IPGeoResponse
	NextSunrise *time.Time
	ISS         ISSMarker
	Anim        *AnimationState
}

type ControlState struct {
	Paused      bool
	Meteors     bool
	Aurora      bool
	Trail       bool
	Pulses      bool
	Scanlines   bool
	SolarSystem bool
	Quit        bool
}

type ZoomState struct {
	Active   bool    // currently in solar system view
	Progress float64 // 0.0 = globe, 1.0 = solar system
}

type Star struct {
	X       int
	Y       int
	Speed   int
	Twinkle int
}

type Meteor struct {
	X    int
	Y    int
	DX   int
	DY   int
	Life int
}

type TrailPoint struct {
	Lat float64
	Lon float64
	Age int
}

type AnimationState struct {
	Stars      []Star
	Meteors    []Meteor
	ISSTrail   []TrailPoint
	BootFrames int
}

var Cities = []City{
	{Name: "New York", Country: "USA", Lat: 40.7128, Lon: -74.0060, Greeting: "Hello"},
	{Name: "Mexico City", Country: "Mexico", Lat: 19.4326, Lon: -99.1332, Greeting: "Hola"},
	{Name: "Rio", Country: "Brazil", Lat: -22.9068, Lon: -43.1729, Greeting: "Ola"},
	{Name: "London", Country: "UK", Lat: 51.5072, Lon: -0.1276, Greeting: "Hello"},
	{Name: "Lagos", Country: "Nigeria", Lat: 6.5244, Lon: 3.3792, Greeting: "Hello"},
	{Name: "Cairo", Country: "Egypt", Lat: 30.0444, Lon: 31.2357, Greeting: "Marhaban"},
	{Name: "Mumbai", Country: "India", Lat: 19.0760, Lon: 72.8777, Greeting: "Namaste"},
	{Name: "Bangkok", Country: "Thailand", Lat: 13.7563, Lon: 100.5018, Greeting: "Sawasdee"},
	{Name: "Tokyo", Country: "Japan", Lat: 35.6762, Lon: 139.6503, Greeting: "Konnichiwa"},
	{Name: "Sydney", Country: "Australia", Lat: -33.8688, Lon: 151.2093, Greeting: "Gday"},
	{Name: "Auckland", Country: "New Zealand", Lat: -36.8509, Lon: 174.7645, Greeting: "Kia ora"},
	{Name: "Honolulu", Country: "USA", Lat: 21.3069, Lon: -157.8583, Greeting: "Aloha"},
}
