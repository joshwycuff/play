package model

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joshwycuff/play/model/command"
	"github.com/joshwycuff/play/model/stdin"
	"github.com/joshwycuff/play/model/stdout"
	"github.com/joshwycuff/play/util"
)

type Model struct {
	ready         bool
	height        int
	width         int
	focus         int
	command       command.Model
	input         stdin.Model
	output        stdout.Model
	inputVisible  bool
	outputVisible bool
	history       History
	keyMap        KeyMap
}

func New(rootCommand string, rootStdout *string) Model {
	history := NewHistory(rootCommand, rootStdout)
	model := Model{
		ready:         false,
		height:        -1,
		width:         -1,
		focus:         1,
		command:       command.New(),
		input:         stdin.New(rootStdout),
		output:        stdout.New(),
		inputVisible:  true,
		outputVisible: true,
		history:       history,
		keyMap:        GetDefaultKeyMap(),
	}
	model.command.Focus()
	return model
}

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		return m.handleWindowSizeMsg(msg)
	default:
		return m.bubbleDown(msg)
	}
}

func (m *Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if key.Matches(msg, m.keyMap.Exit) {
		return m, tea.Quit
	} else if key.Matches(msg, m.keyMap.Enter) {
		return m, m.handleRun()
	} else if key.Matches(msg, m.keyMap.Tab) {
		return m, m.rotateFocus()
	} else if key.Matches(msg, m.keyMap.NavigateToPreviousSibling) {
		m.navigateHistoryToPreviousSibling()
	} else if key.Matches(msg, m.keyMap.NavigateToNextSibling) {
		m.navigateHistoryToNextSibling()
	} else if key.Matches(msg, m.keyMap.NavigateToParent) {
		m.navigateHistoryToParent()
	} else if key.Matches(msg, m.keyMap.NavigateToLatestChild) {
		m.navigateHistoryToLatestChild()
	} else if key.Matches(msg, m.keyMap.ToggleInputVisibility) {
		m.inputVisible = !m.inputVisible
	} else if key.Matches(msg, m.keyMap.ToggleOutputVisibility) {
		m.outputVisible = !m.outputVisible
	} else {
		m.bubbleDownFocus(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.ready = true
	m.height = msg.Height
	m.width = msg.Width

	return m.bubbleDown(msg)
}

func (m *Model) bubbleDown(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	cmd := m.command.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	cmd = m.input.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	cmd = m.output.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) bubbleDownFocus(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.focus {
	case FocusCommand:
		return m, m.command.Update(msg)
	case FocusInput:
		return m, m.input.Update(msg)
	case FocusOutput:
		return m, m.output.Update(msg)
	}
	return nil, nil
}

func (m *Model) handleRun() tea.Cmd {

	historyEntry, err := run(m.command.GetContent(), m.input.GetContent())
	if err != nil {
		m.output.Failure()
		m.output.SetContent(util.P(err.Error()))
		return nil
	}

	node := m.history.AddSibling(&historyEntry)
	m.updateFromNode(node)

	return nil
}

func (m *Model) updateFromNode(node *Node[HistoryData]) {
	if !node.Equals(m.history.GetCurrent()) {
		m.history.SetCurrent(node)
		m.command.SetContent(node.Data.Command)
		m.input.SetContent(node.Data.Stdin)
		if node.Data.ExitStatus == 0 {
			m.output.Success()
			m.output.SetContent(node.Data.Stdout)
		} else {
			m.output.Failure()
			m.output.SetContent(node.Data.Stderr)
		}
	}
}

func (m *Model) navigateHistoryToParent() {
	m.updateFromNode(m.history.GetParent())
}

func (m *Model) navigateHistoryToLatestChild() {
	child := m.history.GetLatestChild()
	if child.Equals(m.history.current) && len(*m.history.GetCurrent().Data.Stdout) > 0 {
		data := NewHistoryData()
		data.Stdin = m.history.GetCurrent().Data.Stdout
		newChild := m.history.AddChild(data)
		m.updateFromNode(newChild)
	} else if !child.Equals(m.history.current) {
		m.updateFromNode(child)
	}
}

func (m *Model) navigateHistoryToPreviousSibling() {
	m.updateFromNode(m.history.GetPrevSibling())
}

func (m *Model) navigateHistoryToNextSibling() {
	m.updateFromNode(m.history.GetNextSibling())
}

func (m Model) View() string {
	if !m.ready {
		return "..."
	}
	var panes []string
	if m.inputVisible {
		panes = append(panes, m.input.View())
	}
	if m.outputVisible {
		panes = append(panes, m.output.View())
	}
	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.command.View(),
		lipgloss.JoinHorizontal(lipgloss.Left, panes...),
	)
}

func (m *Model) rotateFocus() tea.Cmd {
	switch m.focus {
	case FocusCommand:
		m.command.Unfocus()
		m.input.Focus()
		m.focus = FocusInput
	case FocusInput:
		m.input.Unfocus()
		m.output.Focus()
		m.focus = FocusOutput
	case FocusOutput:
		m.output.Unfocus()
		m.command.Focus()
		m.focus = FocusCommand
	}
	return nil
}
