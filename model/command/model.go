package command

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joshwycuff/play/model/common"
)

type Model struct {
	container lipgloss.Style
	textInput textinput.Model
}

func New() Model {
	container := common.GetRoundedBorder()
	textInput := textinput.New()
	textInput.Focus()
	textInput.Placeholder = "Type a command..."
	textInput.Prompt = ""
	return Model{container: container, textInput: textInput}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSizeMsg(msg)
	default:
		return m.bubbleDown(msg)
	}
}

func (m *Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) tea.Cmd {
	m.textInput.Width = msg.Width - 3
	return nil
}

func (m *Model) bubbleDown(msg tea.Msg) tea.Cmd {
	textInput, cmd := m.textInput.Update(msg)
	m.textInput = textInput
	return cmd
}

func (m *Model) View() string {
	return m.container.Render(m.textInput.View())
}

func (m *Model) Content() string {
	return m.textInput.Value()
}

func (m *Model) Focus() {
	m.container = m.container.BorderForeground(lipgloss.Color(common.BLUE))
}

func (m *Model) Unfocus() {
	m.container = m.container.BorderForeground(lipgloss.Color(common.GREY))
}
