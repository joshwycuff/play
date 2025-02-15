package common

import "github.com/charmbracelet/lipgloss"

func GetRoundedBorder() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(CatppuccinMocha.Inactive)
}
