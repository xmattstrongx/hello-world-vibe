package demo

// Theme defines a retro CRT color scheme using ANSI escape codes.
// Monochrome CRT themes use three intensity levels of one phosphor color.
type Theme struct {
	Name   string
	Bright string // high-intensity elements (ISS, cities, aurora highlights)
	Normal string // standard elements (globe terrain, text)
	Dim    string // low-intensity elements (dark terrain, stars, background)
	Reset  string // ANSI reset sequence
}

var themeOrder = []string{"default", "amber", "green", "blue"}

var themes = map[string]Theme{
	"default": {Name: "default"},
	"amber": {
		Name:   "amber",
		Bright: "\033[38;5;214m",
		Normal: "\033[38;5;208m",
		Dim:    "\033[38;5;130m",
		Reset:  "\033[0m",
	},
	"green": {
		Name:   "green",
		Bright: "\033[38;5;46m",
		Normal: "\033[38;5;34m",
		Dim:    "\033[38;5;22m",
		Reset:  "\033[0m",
	},
	"blue": {
		Name:   "blue",
		Bright: "\033[38;5;75m",
		Normal: "\033[38;5;69m",
		Dim:    "\033[38;5;24m",
		Reset:  "\033[0m",
	},
}

// GetTheme returns the theme for a given name, falling back to default.
func GetTheme(name string) Theme {
	if t, ok := themes[name]; ok {
		return t
	}
	return themes["default"]
}

// NextTheme cycles to the next theme in order.
func NextTheme(current string) string {
	for i, name := range themeOrder {
		if name == current {
			return themeOrder[(i+1)%len(themeOrder)]
		}
	}
	return themeOrder[0]
}

// ValidTheme returns true if the name is a recognized theme.
func ValidTheme(name string) bool {
	_, ok := themes[name]
	return ok
}

// runeIntensity classifies a grid rune into bright/normal/dim for theming.
func runeIntensity(ch rune) int {
	switch ch {
	case '@', 'o', '#', '%', '^', '~':
		return 2 // bright
	case '*', '+', '=', '-', ':', ';', ',', '\\':
		return 1 // normal
	default:
		return 0 // dim
	}
}

// colorizeGridLine converts a grid row into a themed string.
// For the default theme (no color codes), it returns the plain string.
func colorizeGridLine(row []rune, t Theme) string {
	if t.Reset == "" {
		return string(row)
	}
	levels := [3]string{t.Dim, t.Normal, t.Bright}
	buf := make([]byte, 0, len(row)*8)
	prev := -1
	for _, ch := range row {
		if ch == ' ' {
			if prev != -1 {
				buf = append(buf, t.Reset...)
				prev = -1
			}
			buf = append(buf, ' ')
			continue
		}
		lvl := runeIntensity(ch)
		if lvl != prev {
			buf = append(buf, levels[lvl]...)
			prev = lvl
		}
		buf = append(buf, string(ch)...)
	}
	if prev != -1 {
		buf = append(buf, t.Reset...)
	}
	return string(buf)
}

// colorizeText wraps a plain text string in the theme's normal color.
func colorizeText(s string, t Theme) string {
	if t.Reset == "" {
		return s
	}
	return t.Normal + s + t.Reset
}

// colorizeBright wraps text in the theme's bright color.
func colorizeBright(s string, t Theme) string {
	if t.Reset == "" {
		return s
	}
	return t.Bright + s + t.Reset
}
