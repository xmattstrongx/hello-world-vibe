package demo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const DefaultInterval = 700 * time.Millisecond

type Config struct {
	Live     bool
	Lang     string
	View     string
	MaxASCII bool
	Frames   int
	Interval time.Duration
}

type Runner struct {
	Out    io.Writer
	Now    func() time.Time
	Sleep  func(time.Duration)
	Client *http.Client
}

func Run(out io.Writer, cfg Config) error {
	if out == nil {
		return errors.New("output writer is nil")
	}
	if cfg.Lang == "" {
		cfg.Lang = "auto"
	}
	if cfg.Lang != "auto" && cfg.Lang != "en" {
		return fmt.Errorf("unsupported --lang value %q (use auto or en)", cfg.Lang)
	}
	if cfg.View == "" {
		cfg.View = "auto"
	}
	if cfg.View != "auto" && cfg.View != "compact" && cfg.View != "cinematic" {
		return fmt.Errorf("unsupported --view value %q (use auto, compact, or cinematic)", cfg.View)
	}
	if cfg.Interval <= 0 {
		cfg.Interval = DefaultInterval
	}

	r := Runner{
		Out:   out,
		Now:   time.Now,
		Sleep: time.Sleep,
		Client: &http.Client{
			Timeout: 4 * time.Second,
		},
	}
	return r.Live(context.Background(), cfg)
}

func (r Runner) Live(ctx context.Context, cfg Config) error {
	if r.Out == nil {
		return errors.New("runner output is nil")
	}
	if r.Now == nil {
		r.Now = time.Now
	}
	if r.Sleep == nil {
		r.Sleep = time.Sleep
	}
	if r.Client == nil {
		r.Client = &http.Client{Timeout: 4 * time.Second}
	}

	if os.Getenv("TERM") == "dumb" {
		fmt.Fprintln(r.Out, "warning: terminal may not support animation")
	}
	restoreUI := beginTerminalUI(r.Out)
	defer restoreUI()

	local, _ := FetchLocalGeo(ctx, r.Client)
	next, _ := NextSunrise(ctx, r.Client, local, r.Now())

	var iss ISSMarker
	lastISSFetch := time.Time{}
	termW, termH := terminalSize(r.Out)
	initialView := resolveViewMode(termW, termH, cfg.View)
	width, height, _ := globeViewport(termW, termH, initialView)
	anim := NewAnimationState(width, height, r.Now().UnixNano())
	controls := defaultControls()
	input, restoreInput, rawMode := startInputListenerWithCleanup(os.Stdin)
	defer restoreInput()
	maxASCII := cfg.MaxASCII || os.Getenv("TERM_PROGRAM") == "iTerm.app"

	frame := 0
	for {
		for {
			select {
			case cmd, ok := <-input:
				if !ok {
					input = nil
					break
				}
				applyCommand(&controls, cmd)
			default:
				goto commandsDrained
			}
		}
	commandsDrained:
		if controls.Quit {
			return nil
		}

		now := r.Now().UTC()
		termW, termH := terminalSize(r.Out)
		view := resolveViewMode(termW, termH, cfg.View)
		width, height, _ := globeViewport(termW, termH, view)
		ResizeAnimationState(&anim, width, height)

		if !controls.Paused {
			AdvanceAnimations(&anim, width, height, frame)
		}
		if !controls.Paused && now.Sub(lastISSFetch) > 10*time.Second {
			if lat, lon, ok := FetchISS(ctx, r.Client); ok {
				iss = ISSMarker{Lat: lat, Lon: lon, OK: true}
			}
			lastISSFetch = now
		}
		if !controls.Paused && iss.OK {
			AddISSTrail(&anim, iss.Lat, iss.Lon)
		}

		clearScreen(r.Out)
		frameData := FrameData{
			Now:         now,
			Lang:        cfg.Lang,
			ViewMode:    view,
			MaxASCII:    maxASCII,
			Frame:       frame,
			RawInput:    rawMode,
			Controls:    controls,
			Local:       local,
			NextSunrise: next,
			ISS:         iss,
			Anim:        &anim,
		}
		fmt.Fprint(r.Out, RenderFrame(frameData, termW, termH))

		if !controls.Paused {
			frame++
		}
		if cfg.Frames > 0 && frame >= cfg.Frames {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		r.Sleep(cfg.Interval)
	}
}

func clearScreen(out io.Writer) {
	fmt.Fprint(out, "\033[H")
}

func RenderFrame(data FrameData, termW, termH int) string {
	view := data.ViewMode
	if view == "" {
		view = resolveViewMode(termW, termH, "auto")
	}
	width, height, leftPad := globeViewport(termW, termH, view)
	decl := SolarDeclination(data.Now)
	subLon := SubsolarLongitude(data.Now)

	grid := makeGrid(width, height)
	mask := renderGlobe(grid, subLon, decl, data.Frame, view, data.MaxASCII)
	if data.Anim != nil {
		drawStarfield(grid, data.Anim.Stars, data.Frame, mask)
		if data.Controls.Meteors {
			drawMeteors(grid, data.Anim.Meteors, mask)
		}
		if data.Controls.Aurora {
			drawAurora(grid, mask, subLon, decl, data.Frame)
		}
	}
	hot := DaylightCities(data.Lang, decl, subLon)
	if data.Anim != nil {
		if data.Controls.Pulses {
			drawCityPulses(grid, hot, subLon, data.Frame)
		}
		if data.Controls.Trail {
			drawISSTrail(grid, data.Anim.ISSTrail, subLon)
		}
	}
	plotCities(grid, hot, subLon)
	if data.ISS.OK {
		plotISS(grid, data.ISS, subLon)
	}
	if data.Anim != nil && data.Controls.Scanlines {
		applyScanlines(grid, data.Frame)
	}

	infoLines := makeInfoLines(data, hot, view)
	maxInfo := termH - (height + 4)
	if maxInfo < 3 {
		maxInfo = 3
	}
	if view == "compact" && maxInfo > 5 {
		maxInfo = 5
	}

	pad := strings.Repeat(" ", leftPad)
	lines := make([]string, 0, termH)
	quality := "std"
	if data.MaxASCII {
		quality = "max"
	}
	header := fitLine(fmt.Sprintf("HELLO, WORLD FROM SPACE | UTC %s | VIEW %s | ASCII %s", data.Now.Format("15:04:05"), view, quality), termW)
	lines = append(lines, header)
	lines = append(lines, pad+strings.Repeat("=", width))
	for _, row := range grid {
		lines = append(lines, pad+string(row))
	}
	lines = append(lines, pad+strings.Repeat("=", width))
	printed := 0
outer:
	for _, ln := range infoLines {
		for _, wrapped := range wrapLine(ln, termW) {
			if printed >= maxInfo {
				break outer
			}
			lines = append(lines, fitLine(wrapped, termW))
			printed++
		}
	}
	return renderFixedFrame(lines, termW, termH)
}

func makeInfoLines(data FrameData, hot []City, view string) []string {
	lines := []string{}
	if len(hot) == 0 {
		lines = append(lines, "Daylight hellos: none in sample set")
	} else {
		names := make([]string, 0, 5)
		limit := 5
		if len(hot) < limit {
			limit = len(hot)
		}
		for i := 0; i < limit; i++ {
			c := hot[i]
			if view == "compact" {
				names = append(names, c.Name)
			} else {
				names = append(names, fmt.Sprintf("%s(%s)", c.Name, c.Greeting))
			}
		}
		lines = append(lines, "Daylight hellos: "+strings.Join(names, "  |  "))
	}

	if data.Local != nil {
		if view == "compact" {
			lines = append(lines, fmt.Sprintf("Location: %s", data.Local.City))
		} else {
			lines = append(lines, fmt.Sprintf("Your location: %s (%.2f, %.2f)", data.Local.City, data.Local.Lat, data.Local.Lon))
		}
	}
	if data.NextSunrise != nil {
		remaining := data.NextSunrise.Sub(data.Now)
		if remaining < 0 {
			remaining = 0
		}
		lines = append(lines, fmt.Sprintf("Next sunrise: %s", FormatDuration(remaining)))
	}
	if data.ISS.OK {
		lines = append(lines, fmt.Sprintf("ISS marker @ lat %.2f lon %.2f", data.ISS.Lat, data.ISS.Lon))
	}
	if view == "cinematic" && data.Anim != nil {
		lines = append(lines, fmt.Sprintf("Effects: stars=%d meteors=%d trail=%d", len(data.Anim.Stars), len(data.Anim.Meteors), len(data.Anim.ISSTrail)))
	}
	lines = append(lines, fmt.Sprintf("Toggles: pause=%t m=%t a=%t t=%t c=%t s=%t", data.Controls.Paused, data.Controls.Meteors, data.Controls.Aurora, data.Controls.Trail, data.Controls.Pulses, data.Controls.Scanlines))
	if data.RawInput {
		if view == "compact" {
			lines = append(lines, "Controls: [space]/p pause | m/a/t/c/s toggle | q quit")
		} else {
			lines = append(lines, "Controls: [space]/p pause | m meteors | a aurora | t trail | c pulses | s scanlines | q quit")
		}
	} else {
		if view == "compact" {
			lines = append(lines, "Controls: p,m,a,t,c,s,q (type + Enter)")
		} else {
			lines = append(lines, "Controls: p(space) pause | m meteors | a aurora | t trail | c pulses | s scanlines | q quit (type + Enter)")
		}
	}

	lines = append(lines, "Language mode: "+data.Lang)
	if view == "cinematic" && data.Anim != nil && data.Frame < data.Anim.BootFrames {
		lines = append(lines, fmt.Sprintf("Boot sequence: T-%02d [scanline init]", data.Anim.BootFrames-data.Frame))
	}
	lines = append(lines, "Ctrl+C to stop")
	return lines
}

func fitLine(line string, width int) string {
	if width <= 0 {
		return ""
	}
	r := []rune(line)
	if len(r) <= width {
		return line
	}
	if width <= 1 {
		return string(r[:1])
	}
	return string(r[:width-1]) + "…"
}

func wrapLine(s string, width int) []string {
	if width < 12 {
		return []string{fitLine(s, width)}
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}
	out := make([]string, 0, 2)
	cur := words[0]
	for _, w := range words[1:] {
		next := cur + " " + w
		if len([]rune(next)) <= width {
			cur = next
			continue
		}
		out = append(out, cur)
		cur = w
	}
	out = append(out, cur)
	return out
}

func renderFixedFrame(lines []string, width, height int) string {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	var b strings.Builder
	for i := 0; i < height; i++ {
		line := ""
		if i < len(lines) {
			line = fitLine(lines[i], width)
		}
		b.WriteString(padRight(line, width))
		if i != height-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func padRight(s string, width int) string {
	r := []rune(s)
	if len(r) >= width {
		return string(r[:width])
	}
	return s + strings.Repeat(" ", width-len(r))
}

func beginTerminalUI(out io.Writer) func() {
	fmt.Fprint(out, "\033[?1049h\033[?25l\033[?7l\033[H")
	return func() {
		fmt.Fprint(out, "\033[?7h\033[?25h\033[?1049l")
	}
}
