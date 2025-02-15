package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joshwycuff/play/model/command"
	"github.com/joshwycuff/play/model/stdin"
	"github.com/joshwycuff/play/model/stdout"
)

type Model struct {
	ready   bool
	height  int
	width   int
	focus   int
	command command.Model
	input   stdin.Model
	output  stdout.Model
}

func New(inputContent string) Model {
	cmd := command.New()
	cmd.Focus()
	return Model{
		ready:   false,
		height:  -1,
		width:   -1,
		focus:   1,
		command: cmd,
		input:   stdin.New(inputContent),
		output:  stdout.New(),
	}
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

	key := msg.String()
	if key == "ctrl+c" || key == "esc" {
		return m, tea.Quit
	} else if key == "enter" {
		return m, m.handleRun()
	} else if key == "tab" {
		return m, m.rotateFocus()
	} else if m.focus == FocusCommand && key == "up" {
		m.handleHistoryPrev()
	} else if m.focus == FocusCommand && key == "down" {
		m.handleHistoryNext()
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

	historyEntry, err := run(m.command.Content(), m.input.Content())
	if err != nil {
		m.output.Failure()
		m.output.SetContent(err.Error())
		return nil
	}

	if historyEntry.result.ExitStatus == 0 {
		m.output.Success()
		m.output.SetContent(historyEntry.result.Stdout)
	} else {
		m.output.Failure()
		m.output.SetContent(historyEntry.result.Stderr)
	}

	PushHistoryEntry(historyEntry)

	return nil
}

func (m *Model) handleHistoryPrev() {
	historyEntry := GetPreviousHistoryEntry()
	m.command.SetContent(historyEntry.command)
	if historyEntry.result.ExitStatus == 0 {
		m.output.Success()
		m.output.SetContent(historyEntry.result.Stdout)
	} else {
		m.output.Failure()
		m.output.SetContent(historyEntry.result.Stderr)
	}
}

func (m *Model) handleHistoryNext() {
	historyEntry := GetNextHistoryEntry()
	m.command.SetContent(historyEntry.command)
	if historyEntry.result.ExitStatus == 0 {
		m.output.Success()
		m.output.SetContent(historyEntry.result.Stdout)
	} else {
		m.output.Failure()
		m.output.SetContent(historyEntry.result.Stderr)
	}
}

func (m Model) View() string {
	if !m.ready {
		return "..."
	}
	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.command.View(),
		lipgloss.JoinHorizontal(lipgloss.Left, m.input.View(), m.output.View()),
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
