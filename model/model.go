package model

import (
	"bytes"
	"os/exec"

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

	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit
	case "enter":
		return m, m.handleRun()
	case "tab":
		return m, m.rotateFocus()
	default:
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
	case 1:
		return m, m.command.Update(msg)
	case 2:
		return m, m.input.Update(msg)
	case 3:
		return m, m.output.Update(msg)
	}
	return nil, nil
}

func (m *Model) handleRun() tea.Cmd {
	cmd := exec.Command("sh", "-c", m.command.Content())

	// Get the command's stdin pipe
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil
	}

	// Provide input to the command
	go func() {
		defer stdin.Close()
		stdin.Write([]byte(m.input.Content()))
	}()

	// Capture output
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command
	if err := cmd.Run(); err != nil {
		m.output.Failure()
		m.output.SetContent(err.Error())
		return nil
	}

	m.output.Success()
	m.output.SetContent(out.String())
	return nil
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
	case 1:
		m.command.Unfocus()
		m.input.Focus()
		m.focus = 2
	case 2:
		m.input.Unfocus()
		m.output.Focus()
		m.focus = 3
	case 3:
		m.output.Unfocus()
		m.command.Focus()
		m.focus = 1
	}
	return nil
}
