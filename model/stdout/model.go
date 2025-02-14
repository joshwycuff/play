package stdout

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joshwycuff/play/model/common"
)

type Model struct {
	container lipgloss.Style
	viewPort  viewport.Model
	focus     bool
	failure   bool
}

func New() Model {
	container := common.GetRoundedBorder()
	viewPort := viewport.New(-1, -1)
	return Model{container: container, viewPort: viewPort, focus: false, failure: false}
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
	color := common.GREY
	if m.failure {
		color = common.RED
	} else if m.focus {
		color = common.BLUE
	}
	return m.container.BorderForeground(color).Render(m.viewPort.View())
}

func (m *Model) SetContent(content string) {
	m.viewPort.SetContent(content)
}

func (m *Model) Focus() {
	m.focus = true
}

func (m *Model) Unfocus() {
	m.focus = false
}

func (m *Model) Success() {
	m.failure = false
}

func (m *Model) Failure() {
	m.failure = true
}
