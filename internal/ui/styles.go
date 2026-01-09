package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Note name highlighting (cyan/bold)
	NoteName = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true)

	// Success message
	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color("82"))

	// Error message
	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196"))

	// Interactive input box
	InputBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1)

	// Prompt/hint text
	Hint = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true)
)
