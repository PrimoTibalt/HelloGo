package main

import (
	"fmt"
	"os"
	retriever "primotibalt/checkTests/questionsRetriever"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const (
	paddingToLeft           = 4
	defaultTaPlaceholder    = "Напиши ответ на вопрос"
	delimeterQuestionAnswer = "/!/"
)

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
					m = initializeModel(slices.Concat(m.SuccessQuestions, m.FailedQuestions))
					return m, tea.Batch(taCmd, vpCmd, tea.ClearScreen)
				case "c":
					fmt.Println("Ты не имеешь ничего для выбора")
					return m, tea.Quit
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
			"Нажми r(reset)/c(choose) для продолжения",
			delimeter)
	}
}

func main() {
	for {
		questions := []Question{}
		topics := retriever.RetrieveTopicToPathMap()
		options := make([]huh.Option[string], len(topics))
		var i int
		for topic, path := range topics {
			options[i] = huh.NewOption(topic, path)
			i++
		}
		var selectedPaths []string
		err := huh.NewMultiSelect[string]().
			Title("Выбери топик(и) для теста:").
			Options(options...).
			Value(&selectedPaths).
			Run()
		if err != nil {
			if err == huh.ErrUserAborted {
				os.Exit(0)
			}
		}

		if len(selectedPaths) == 0 {
			continue
		}

		for _, path := range selectedPaths {
			qaPairs := retriever.TopicQuestions(path)
			for _, pair := range qaPairs {
				if !strings.Contains(pair, delimeterQuestionAnswer) {
					continue
				}

				qa := strings.Split(pair, delimeterQuestionAnswer)
				question := qa[0]
				answer := qa[1]
				questions = append(questions, Question{answer, question})
			}
		}

		program := tea.NewProgram(initializeModel(questions))
		_, err = program.Run()
		if err != nil {
			panic(err)
		}
	}
}
