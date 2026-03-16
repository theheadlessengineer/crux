package tui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Theme holds the semantic colour palette for the TUI, mirroring llmfit's ThemeColors.
type Theme struct {
	Name string

	// General
	BG     lipgloss.Color
	FG     lipgloss.Color
	Muted  lipgloss.Color
	Border lipgloss.Color
	Title  lipgloss.Color

	// Interaction
	Accent   lipgloss.Color
	Cursor   lipgloss.Color
	Selected lipgloss.Color

	// Status
	Good    lipgloss.Color
	Warning lipgloss.Color
	Error   lipgloss.Color

	// Status bar (inverted)
	StatusBG lipgloss.Color
	StatusFG lipgloss.Color
}

// Styles derives lipgloss.Style values from the theme's colour tokens.
type Styles struct {
	Border    lipgloss.Style
	Title     lipgloss.Style
	Prompt    lipgloss.Style
	Cursor    lipgloss.Style
	Selected  lipgloss.Style
	Muted     lipgloss.Style
	Error     lipgloss.Style
	Success   lipgloss.Style
	Accent    lipgloss.Style
	Help      lipgloss.Style
	Panel     lipgloss.Style
	StatusBar lipgloss.Style
}

// BuildStyles derives all styles from the theme.
func (t *Theme) BuildStyles() Styles {
	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border)
	if t.BG != "" {
		panel = panel.Background(t.BG)
	}

	statusBar := lipgloss.NewStyle().
		Background(t.StatusBG).
		Foreground(t.StatusFG).
		Bold(true)

	title := lipgloss.NewStyle().Bold(true).Foreground(t.Title)
	if t.BG != "" {
		title = title.Background(t.BG)
	}

	return Styles{
		Border:    lipgloss.NewStyle().Foreground(t.Border),
		Title:     title,
		Prompt:    lipgloss.NewStyle().Bold(true).Foreground(t.FG),
		Cursor:    lipgloss.NewStyle().Foreground(t.Cursor).Bold(true),
		Selected:  lipgloss.NewStyle().Foreground(t.Selected).Bold(true),
		Muted:     lipgloss.NewStyle().Foreground(t.Muted),
		Error:     lipgloss.NewStyle().Foreground(t.Error),
		Success:   lipgloss.NewStyle().Foreground(t.Good),
		Accent:    lipgloss.NewStyle().Foreground(t.Accent),
		Help:      lipgloss.NewStyle().Foreground(t.Muted),
		Panel:     panel,
		StatusBar: statusBar,
	}
}

var themes = []Theme{
	{
		Name: "default",
		BG:   "", FG: lipgloss.Color("255"), Muted: lipgloss.Color("240"),
		Border: lipgloss.Color("240"), Title: lipgloss.Color("82"),
		Accent: lipgloss.Color("75"), Cursor: lipgloss.Color("75"), Selected: lipgloss.Color("75"),
		Good: lipgloss.Color("82"), Warning: lipgloss.Color("220"), Error: lipgloss.Color("196"),
		StatusBG: lipgloss.Color("82"), StatusFG: lipgloss.Color("0"),
	},
	{
		Name: "dracula",
		BG:   lipgloss.Color("#282a36"), FG: lipgloss.Color("#f8f8f2"), Muted: lipgloss.Color("#6272a4"),
		Border: lipgloss.Color("#44475a"), Title: lipgloss.Color("#50fa7b"),
		Accent: lipgloss.Color("#8be9fd"), Cursor: lipgloss.Color("#bd93f9"), Selected: lipgloss.Color("#50fa7b"),
		Good: lipgloss.Color("#50fa7b"), Warning: lipgloss.Color("#f1fa8c"), Error: lipgloss.Color("#ff5555"),
		StatusBG: lipgloss.Color("#bd93f9"), StatusFG: lipgloss.Color("#282a36"),
	},
	{
		Name: "solarized",
		BG:   lipgloss.Color("#002b36"), FG: lipgloss.Color("#839496"), Muted: lipgloss.Color("#586e75"),
		Border: lipgloss.Color("#586e75"), Title: lipgloss.Color("#859900"),
		Accent: lipgloss.Color("#268bd2"), Cursor: lipgloss.Color("#268bd2"), Selected: lipgloss.Color("#859900"),
		Good: lipgloss.Color("#859900"), Warning: lipgloss.Color("#b58900"), Error: lipgloss.Color("#dc322f"),
		StatusBG: lipgloss.Color("#268bd2"), StatusFG: lipgloss.Color("#fdf6e3"),
	},
	{
		Name: "nord",
		BG:   lipgloss.Color("#2e3440"), FG: lipgloss.Color("#d8dee9"), Muted: lipgloss.Color("#4c566a"),
		Border: lipgloss.Color("#434c5e"), Title: lipgloss.Color("#a3be8c"),
		Accent: lipgloss.Color("#88c0d0"), Cursor: lipgloss.Color("#88c0d0"), Selected: lipgloss.Color("#a3be8c"),
		Good: lipgloss.Color("#a3be8c"), Warning: lipgloss.Color("#ebcb8b"), Error: lipgloss.Color("#bf616a"),
		StatusBG: lipgloss.Color("#81a1c1"), StatusFG: lipgloss.Color("#2e3440"),
	},
	{
		Name: "monokai",
		BG:   lipgloss.Color("#272822"), FG: lipgloss.Color("#f8f8f2"), Muted: lipgloss.Color("#75715e"),
		Border: lipgloss.Color("#49483e"), Title: lipgloss.Color("#a6e22e"),
		Accent: lipgloss.Color("#66d9e8"), Cursor: lipgloss.Color("#66d9e8"), Selected: lipgloss.Color("#a6e22e"),
		Good: lipgloss.Color("#a6e22e"), Warning: lipgloss.Color("#e6db74"), Error: lipgloss.Color("#f92672"),
		StatusBG: lipgloss.Color("#f92672"), StatusFG: lipgloss.Color("#272822"),
	},
	{
		Name: "gruvbox",
		BG:   lipgloss.Color("#282828"), FG: lipgloss.Color("#ebdbb2"), Muted: lipgloss.Color("#928374"),
		Border: lipgloss.Color("#504945"), Title: lipgloss.Color("#b8bb26"),
		Accent: lipgloss.Color("#83a598"), Cursor: lipgloss.Color("#83a598"), Selected: lipgloss.Color("#b8bb26"),
		Good: lipgloss.Color("#b8bb26"), Warning: lipgloss.Color("#fabd2f"), Error: lipgloss.Color("#fb4934"),
		StatusBG: lipgloss.Color("#d79921"), StatusFG: lipgloss.Color("#282828"),
	},
}

func themePreferencePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".crux", "theme")
}

// LoadThemePreference reads the saved theme name; returns the default theme if absent.
func LoadThemePreference() Theme {
	p := themePreferencePath()
	if p == "" {
		return themes[0]
	}
	// #nosec G304 — path is constructed from os.UserHomeDir(), not user input
	data, err := os.ReadFile(p)
	if err != nil {
		return themes[0]
	}
	name := strings.TrimSpace(string(data))
	for i := range themes {
		if themes[i].Name == name {
			return themes[i]
		}
	}
	return themes[0]
}

// SaveThemePreference persists the theme name to ~/.crux/theme.
func SaveThemePreference(t *Theme) {
	p := themePreferencePath()
	if p == "" {
		return
	}
	_ = os.MkdirAll(filepath.Dir(p), 0o750)
	_ = os.WriteFile(p, []byte(t.Name+"\n"), 0o600)
}

// NextTheme cycles to the next theme in the list.
func NextTheme(current *Theme) Theme {
	for i := range themes {
		if themes[i].Name == current.Name {
			return themes[(i+1)%len(themes)]
		}
	}
	return themes[0]
}

// ThemeNames returns all available theme names.
func ThemeNames() []string {
	names := make([]string, len(themes))
	for i := range themes {
		names[i] = themes[i].Name
	}
	return names
}
