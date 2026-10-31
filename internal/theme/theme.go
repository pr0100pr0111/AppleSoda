package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

type Theme struct {
	Name           string
	Primary        color.Color
	Secondary      color.Color
	Accent         color.Color
	Background     color.Color
	Foreground     color.Color
	Muted          color.Color
	Border         color.Color
	NowPlaying     color.Color
	ProgressFilled color.Color
	ProgressEmpty  color.Color
	VizGradient    []color.Color
}

var themes = map[string]Theme{
	"apple":      apple(),
	"catppuccin": catppuccin(),
	"gruvbox":    gruvbox(),
	"nord":       nord(),
	"mono":       mono(),
}

var themeOrder = []string{"apple", "catppuccin", "gruvbox", "nord", "mono"}

func Get(name string) Theme {
	if t, ok := themes[name]; ok {
		return t
	}
	return themes["apple"]
}

func Names() []string {
	result := make([]string, len(themeOrder))
	copy(result, themeOrder)
	return result
}

func NextTheme(current string) string {
	for i, name := range themeOrder {
		if name == current {
			return themeOrder[(i+1)%len(themeOrder)]
		}
	}
	return themeOrder[0]
}

func apple() Theme {
	return Theme{
		Name:           "apple",
		Primary:        lipgloss.Color("#fc3c44"),
		Secondary:      lipgloss.Color("#fa57c1"),
		Accent:         lipgloss.Color("#fc3c44"),
		Background:     lipgloss.Color("#1c1c1e"),
		Foreground:     lipgloss.Color("#f5f5f7"),
		Muted:          lipgloss.Color("#86868b"),
		Border:         lipgloss.Color("#3a3a3c"),
		NowPlaying:     lipgloss.Color("#fc3c44"),
		ProgressFilled: lipgloss.Color("#fc3c44"),
		ProgressEmpty:  lipgloss.Color("#3a3a3c"),
		VizGradient: []color.Color{
			lipgloss.Color("#fc3c44"),
			lipgloss.Color("#fa57c1"),
			lipgloss.Color("#af52de"),
		},
	}
}

func catppuccin() Theme {
	return Theme{
		Name:           "catppuccin",
		Primary:        lipgloss.Color("#cba6f7"),
		Secondary:      lipgloss.Color("#f5c2e7"),
		Accent:         lipgloss.Color("#89b4fa"),
		Background:     lipgloss.Color("#1e1e2e"),
		Foreground:     lipgloss.Color("#cdd6f4"),
		Muted:          lipgloss.Color("#6c7086"),
		Border:         lipgloss.Color("#313244"),
		NowPlaying:     lipgloss.Color("#cba6f7"),
		ProgressFilled: lipgloss.Color("#89b4fa"),
		ProgressEmpty:  lipgloss.Color("#313244"),
		VizGradient: []color.Color{
			lipgloss.Color("#89b4fa"),
			lipgloss.Color("#cba6f7"),
			lipgloss.Color("#f5c2e7"),
		},
	}
}

func gruvbox() Theme {
	return Theme{
		Name:           "gruvbox",
		Primary:        lipgloss.Color("#fe8019"),
		Secondary:      lipgloss.Color("#fabd2f"),
		Accent:         lipgloss.Color("#b8bb26"),
		Background:     lipgloss.Color("#282828"),
		Foreground:     lipgloss.Color("#ebdbb2"),
		Muted:          lipgloss.Color("#928374"),
		Border:         lipgloss.Color("#3c3836"),
		NowPlaying:     lipgloss.Color("#fe8019"),
		ProgressFilled: lipgloss.Color("#b8bb26"),
		ProgressEmpty:  lipgloss.Color("#3c3836"),
		VizGradient: []color.Color{
			lipgloss.Color("#b8bb26"),
			lipgloss.Color("#fabd2f"),
			lipgloss.Color("#fe8019"),
		},
	}
}

func nord() Theme {
	return Theme{
		Name:           "nord",
		Primary:        lipgloss.Color("#88c0d0"),
		Secondary:      lipgloss.Color("#81a1c1"),
		Accent:         lipgloss.Color("#5e81ac"),
		Background:     lipgloss.Color("#2e3440"),
		Foreground:     lipgloss.Color("#eceff4"),
		Muted:          lipgloss.Color("#4c566a"),
		Border:         lipgloss.Color("#3b4252"),
		NowPlaying:     lipgloss.Color("#88c0d0"),
		ProgressFilled: lipgloss.Color("#88c0d0"),
		ProgressEmpty:  lipgloss.Color("#3b4252"),
		VizGradient: []color.Color{
			lipgloss.Color("#5e81ac"),
			lipgloss.Color("#81a1c1"),
			lipgloss.Color("#88c0d0"),
		},
	}
}

func mono() Theme {
	return Theme{
		Name:           "mono",
		Primary:        lipgloss.Color("#ffffff"),
		Secondary:      lipgloss.Color("#cccccc"),
		Accent:         lipgloss.Color("#ffffff"),
		Background:     lipgloss.Color("#000000"),
		Foreground:     lipgloss.Color("#ffffff"),
		Muted:          lipgloss.Color("#666666"),
		Border:         lipgloss.Color("#333333"),
		NowPlaying:     lipgloss.Color("#ffffff"),
		ProgressFilled: lipgloss.Color("#ffffff"),
		ProgressEmpty:  lipgloss.Color("#333333"),
		VizGradient: []color.Color{
			lipgloss.Color("#666666"),
			lipgloss.Color("#999999"),
			lipgloss.Color("#ffffff"),
		},
	}
}
