package main

import (
	"fmt"
	"slices"
	"strings"

	persistence "primotibalt/checkTests/topicpersistence"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func ChooseTopicsForTest(topics []persistence.TopicName) (selectedTopics []persistence.TopicName) {
	options := make([]huh.Option[persistence.TopicName], len(topics))
	for ptr, topic := range topics {
		options[ptr] = huh.NewOption(string(topic), topic)
	}
	err := huh.NewForm(huh.NewGroup(huh.NewMultiSelect[persistence.TopicName]().
		Title("Выбери топик(и) для теста:").
		Options(options...).
		Value(&selectedTopics))).
		WithProgramOptions(tea.WithAltScreen()).
		Run()
	if err != nil {
		if err != huh.ErrUserAborted {
			panic(err)
		} else {
			return []persistence.TopicName{}
		}
	}

	return
}

func RunKnowledgeTest(selectedTopics []persistence.TopicName) {
	questions := []Question{}
	for _, topic := range selectedTopics {
		qaPairs := persistence.TopicQuestions(topic)
		for _, pair := range qaPairs {
			if !strings.Contains(pair, persistence.DelimeterQuestionAnswer) {
				continue
			}

			qa := strings.Split(pair, persistence.DelimeterQuestionAnswer)
			question := qa[0]
			answer := qa[1]
			questions = append(questions, Question{answer, question})
		}
	}

	model, err := initializeModel(questions)
	if err != nil {
		panic(err)
	}

	program := tea.NewProgram(model, tea.WithAltScreen())
	_, err = program.Run()
	if err != nil {
		panic(err)
	}
}

func (m TestCheck) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, tea.ClearScreen)
}

func (m TestCheck) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd       tea.Cmd
		vpCmd       tea.Cmd
		vpFailedCmd tea.Cmd
	)

	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.vpFailed, vpFailedCmd = m.vpFailed.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if len(m.Questions) < 1 {
				switch strings.Trim(m.textarea.Value(), " \n") {
				case "r", "R", "reset":
					m, err := initializeModel(slices.Concat(m.SuccessQuestions, m.FailedQuestions))
					if err != nil {
						panic(err)
					}
					return m, tea.Batch(taCmd, vpCmd, tea.ClearScreen)
				default:
					return m, tea.Quit
				}
			}

			distance := m.computeDistanceBetweenInputAndAnswer()
			m.LastQuestionSuccess = m.IsInputAndAnswerEqual(distance)
			if m.LastQuestionSuccess {
				m.SuccessQuestions = append(m.SuccessQuestions, *m.CurrentQuestion)
			} else {
				m.FailedQuestions = append(m.FailedQuestions, *m.CurrentQuestion)
				m.prepareFailedVpContent(strings.Trim(m.textarea.Value(), " \n\r"))
			}

			m.textarea.Reset()
			m.prepareNextQuestion()
		}
	default:
		if m.textarea.Length() == m.textarea.Width()*m.textarea.Height() {
			currentContent := m.textarea.Value()
			m.textarea.SetHeight(m.textarea.Height() + 1)
			m.textarea.SetValue(currentContent)
		}
	}

	return m, tea.Batch(taCmd, vpCmd, vpFailedCmd)
}

func (m TestCheck) View() string {
	const delimeter = "\n"

	if len(m.Questions) > 0 {
		return lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.JoinVertical(lipgloss.Left, m.viewport.View(), m.textarea.View()),
			m.vpFailed.View())
	} else {
		return fmt.Sprintf("%s%s%s%s",
			lipgloss.JoinHorizontal(lipgloss.Left,
				lipgloss.JoinVertical(lipgloss.Left, m.viewport.View(), m.textarea.View()),
				m.vpFailed.View()),
			delimeter,
			"Нажми r(reset) для продолжения. Любую другую для возвращения к списку.",
			delimeter)
	}
}
