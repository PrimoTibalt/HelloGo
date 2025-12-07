package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Notes struct {
	Notes []Note
}

func (n Notes) View() string {
	builder := strings.Builder{}
	for _, note := range n.Notes {
		builder.WriteString(note.Text)
		builder.WriteString("\t")
		builder.WriteString(note.RegistrationDate.Format(TimeFormat))
		if !note.DueDate.IsZero() {
			builder.WriteString("\t")
			builder.WriteString(note.DueDate.Format(TimeFormat))
		}
	}

	builder.WriteRune('\n')
	return builder.String()
}

func (n Notes) Init() tea.Cmd {
	return nil
}

func (n Notes) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return n, tea.Quit
}

func displayNotes(notes []Note) {
	p := tea.NewProgram(Notes{notes})
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
