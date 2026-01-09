package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true).
			MarginBottom(1)

	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)
)

type inputModel struct {
	textarea textarea.Model
	noteName string
	done     bool
	quit     bool
}

func initialModel(noteName string) inputModel {
	ta := textarea.New()
	ta.Placeholder = "Start typing..."
	ta.Focus()
	ta.CharLimit = 0 // No limit
	ta.SetWidth(60)
	ta.SetHeight(6)

	// Style the textarea
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle.Base = lipgloss.NewStyle()
	ta.BlurredStyle.Base = lipgloss.NewStyle()

	return inputModel{
		textarea: ta,
		noteName: noteName,
	}
}

func (m inputModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlD:
			m.done = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quit = true
			return m, tea.Quit
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	var b strings.Builder

	title := titleStyle.Render(fmt.Sprintf("Adding to %s", m.noteName))
	b.WriteString(title)
	b.WriteString("\n")

	content := containerStyle.Render(m.textarea.View())
	b.WriteString(content)
	b.WriteString("\n")

	help := helpStyle.Render("Ctrl+D to save • Esc to cancel")
	b.WriteString(help)

	return b.String()
}

// RunInteractiveInput launches the interactive textarea and returns the content
func RunInteractiveInput(noteName string) (string, bool) {
	p := tea.NewProgram(initialModel(noteName))
	m, err := p.Run()
	if err != nil {
		return "", false
	}

	model := m.(inputModel)
	if model.quit {
		return "", false
	}

	return model.textarea.Value(), model.done
}
