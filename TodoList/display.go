package main

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Notes struct {
	Notes []Note
	Table table.Model
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder())

func (n Notes) View() string {
	return baseStyle.Render(n.Table.View()) + "\n"
}

func (n Notes) Init() tea.Cmd {
	return nil
}

func (n Notes) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return n, tea.Quit
}

func displayNotes(notes []Note) {
	columns := []table.Column{
		{Title: "Order", Width: 6},
		{Title: "Text", Width: 40},
		{Title: "RegistrationDate", Width: 18},
		{Title: "DueDate", Width: 10},
	}

	rows := []table.Row{}
	for i, l := range notes {
		var dueDate string
		if !l.DueDate.IsZero() {
			dueDate = l.DueDate.Format(TimeFormat)
		} else {
			dueDate = ""
		}
		rows = append(rows,
			table.Row{
				strconv.Itoa(i),
				l.Text,
				l.RegistrationDate.Format(TimeFormat),
				dueDate,
			})
	}

	tm := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(len(notes)+1),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.Foreground(lipgloss.Color("0")).Bold(false)
	tm.SetStyles(s)
	p := tea.NewProgram(Notes{notes, tm})
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
