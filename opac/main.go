package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

var start = &Model{
	Choices: []string{
		"Buy carrots",
		"Buy celery",
		"Buy kohlrabi",
	},
	Selected: make(map[int]struct{}),
}

type Model struct {
	Choices  []string         // items on the to-do list
	Cursor   int              // which to-do list item our cursor is pointing at
	Selected map[int]struct{} // which to-do items are selected
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Choices) {
				m.Cursor++
			}
		case "enter", " ":
			_, ok := m.Selected[m.Cursor]
			if ok {
				delete(m.Selected, m.Cursor)
			} else {
				m.Selected[m.Cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m *Model) View() string {
	s := "What should we buy at the market?\n\n"
	for i, choice := range m.Choices {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.Selected[i]; ok {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	s += "\nPress q to quit.\n"
	return s
}

func main() {
	p := tea.NewProgram(start)
	if err := p.Start(); err != nil {
		log.Fatalf("could not start: %v", err)
	}
}
