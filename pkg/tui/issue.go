package tui

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/patya3/notime/pkg/tui/constants"
)

type mode int

const (
	nav mode = iota
	edit
	create
)

type Model struct {
	list  list.Model
	input textinput.Model
	mode
}

func InitIssue() tea.Model {
	input := textinput.New()
	input.Prompt = "$ "
	input.Placeholder = "Project name..."
	input.CharLimit = 250
	input.Width = 50

	projects, err := constants.IssueRepo.GetAllIssues()
	if err != nil {
		log.Fatal(err)
	}
	items := make([]list.Item, len(projects))
	for i, proj := range projects {
		items[i] = list.Item(proj)
	}

	m := Model{list: list.New(items, list.NewDefaultDelegate(), 0, 0), input: input}
	m.list.Title = "Issues"
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			constants.Keymap.Create,
			constants.Keymap.Rename,
			constants.Keymap.Delete,
			constants.Keymap.Back,
		}
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.input.Focused() {
			if key.Matches(msg, constants.Keymap.Enter) {
				if m.mode.create {
					cmds = cmds.append()
				}

			}

		} else {
			switch {
			case key.Matches(msg, constants.Keymap.Quit):
				return m, tea.Quit
			case key.Matches(msg, constants.Keymap.Create):
				m.input.Focus()
				cmd = textinput.Blink

			default:
				m.list, cmd = m.list.Update(msg)
			}
		}
	case tea.WindowSizeMsg:
		h, v := constants.DocStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.input.Focused() {
		return constants.DocStyle.Render(m.list.View() + "\n" + m.input.View())
	}
	return constants.DocStyle.Render(m.list.View())
}
