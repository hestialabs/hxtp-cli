package ui

import "github.com/charmbracelet/lipgloss"

// Theme defines the 2026 HxTP CLI aesthetic.
type Theme struct {
	Primary    lipgloss.AdaptiveColor
	Secondary  lipgloss.AdaptiveColor
	Accent     lipgloss.AdaptiveColor
	Success    lipgloss.AdaptiveColor
	WarningColor lipgloss.AdaptiveColor
	Error      lipgloss.AdaptiveColor

	BoldHeader lipgloss.Style
	SubHeader  lipgloss.Style
	SuccessMsg lipgloss.Style
	ErrorMsg   lipgloss.Style
	InfoMsg    lipgloss.Style
	WarningMsg lipgloss.Style
	AccentMsg  lipgloss.Style
	CodeBlock  lipgloss.Style
}

// GetTheme returns the adaptive 2026 theme for dark and light terminals.
func GetTheme() *Theme {
	t := &Theme{
		Primary:   lipgloss.AdaptiveColor{Light: "#2D3436", Dark: "#DFE6E9"},
		Secondary: lipgloss.AdaptiveColor{Light: "#636E72", Dark: "#B2BEC3"},
		Accent:    lipgloss.AdaptiveColor{Light: "#0984E3", Dark: "#74B9FF"}, 
		Success:   lipgloss.AdaptiveColor{Light: "#00B894", Dark: "#55E6C1"},
		WarningColor: lipgloss.AdaptiveColor{Light: "#E17055", Dark: "#FAB1A0"},
		Error:     lipgloss.AdaptiveColor{Light: "#D63031", Dark: "#FF7675"},
	}

	t.BoldHeader = lipgloss.NewStyle().
		Foreground(t.Accent).
		Bold(true).
		MarginBottom(1)

	t.SubHeader = lipgloss.NewStyle().
		Foreground(t.Secondary).
		Italic(true)

	t.SuccessMsg = lipgloss.NewStyle().
		Foreground(t.Success).
		Bold(true)

	t.ErrorMsg = lipgloss.NewStyle().
		Foreground(t.Error).
		Bold(true)

	t.WarningMsg = lipgloss.NewStyle().
		Foreground(t.WarningColor).
		Bold(true)

	t.AccentMsg = lipgloss.NewStyle().
		Foreground(t.Accent).
		Bold(true)

	t.InfoMsg = lipgloss.NewStyle().
		Foreground(t.Secondary)

	t.CodeBlock = lipgloss.NewStyle().
		Background(lipgloss.AdaptiveColor{Light: "#F1F2F6", Dark: "#2F3542"}).
		Padding(0, 1).
		Foreground(lipgloss.AdaptiveColor{Light: "#2F3542", Dark: "#F1F2F6"})

	return t
}
