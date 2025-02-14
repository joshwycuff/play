package stdin

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joshwycuff/play/model/common"
)

type Model struct {
	container lipgloss.Style
	content   string
	viewPort  viewport.Model
}

func New(content string) Model {
	container := common.GetRoundedBorder()
	viewPort := viewport.New(-1, -1)
	viewPort.SetContent(content)
	return Model{container: container, content: content, viewPort: viewPort}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowSizeMsg(msg)
	}

	cmd := m.bubbleDown(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (m *Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) {
	m.container = m.container.Height(msg.Height - 5).Width(msg.Width/2 - 2)
	m.viewPort.Height = msg.Height - 5
	m.viewPort.Width = msg.Width/2 - 2
}

func (m *Model) bubbleDown(msg tea.Msg) tea.Cmd {
	viewPort, cmd := m.viewPort.Update(msg)
	m.viewPort = viewPort
	return cmd
}

func (m *Model) View() string {
	return m.container.Render(m.viewPort.View())
}

func (m *Model) Content() string {
	return m.content
}

func (m *Model) Focus() {
	m.container = m.container.BorderForeground(lipgloss.Color(common.BLUE))
}

func (m *Model) Unfocus() {
	m.container = m.container.BorderForeground(lipgloss.Color(common.GREY))
}
