package main

import (
	"fmt"
	"math/rand"
	"slices"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	paddingToLeft               = 4
	defaultTaPlaceholder        = "Напиши ответ на вопрос"
	failedQuestionTaPlaceholder = "Нажми Enter чтобы продолжить"
	delimeterQuestionAnswer     = "/!/"
)

func (m TestCheck) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, tea.ClearScreen)
}

func (m TestCheck) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width - paddingToLeft)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if len(m.Questions) < 1 {
				switch msg.String() {
				case "r", "R":
					m = initializeModel()
					return m, tea.Batch(taCmd, vpCmd)
				case "c":
					fmt.Println("Ты не имеешь ничего для выбора")
					return m, tea.Quit
				}
			}

			if !m.LastQuestionSuccess {
				m.LastQuestionSuccess = true
				m.textarea.Reset()
				m.textarea.Placeholder = defaultTaPlaceholder
				return m, nil
			}

			m.LastQuestionSuccess = m.isInputAndAnswerEqual()
			if m.LastQuestionSuccess {
				m.SuccessQuestions = append(m.SuccessQuestions, *m.CurrentQuestion)
			} else {
				m.textarea.Placeholder = failedQuestionTaPlaceholder
				m.FailedQuestions = append(m.FailedQuestions, *m.CurrentQuestion)
				splitAnswer := ""
				answerIndex := 0
				widthOfAnswer := utf8.RuneCountInString(m.CurrentQuestion.Answer)
				widthOfVp := m.viewport.Width - 2
				for answerIndex+widthOfVp < widthOfAnswer {
					splitAnswer += m.CurrentQuestion.Answer[answerIndex : answerIndex+widthOfVp]
					answerIndex += widthOfVp
				}

				splitAnswer += m.CurrentQuestion.Answer[answerIndex : widthOfAnswer-1]
				m.vpFailed.SetContent(
					fmt.Sprintf("Это неправильный ответ! Правильный ответ такой:\n\033[1m%s\033[0m\n%s",
						m.CurrentQuestion.Answer,
						"Нажмите на любую клавишу для продолжения"))
			}

			m.textarea.Reset()

			questionIndex := slices.Index(m.Questions, *m.CurrentQuestion)
			m.Questions = slices.Delete(m.Questions, questionIndex, questionIndex+1)

			questionsToAskCount := len(m.Questions)
			if questionsToAskCount > 0 {
				selectedQuestionIndex := rand.Intn(questionsToAskCount)
				m.CurrentQuestion = &m.Questions[selectedQuestionIndex]
				m.viewport.SetContent(m.CurrentQuestion.Text)
			} else {
				m.setContentWhenNoMoreQuestions()
			}
		}
	_:
		if m.textarea.Length() == m.textarea.Width()*m.textarea.Height() {
			currentContent := m.textarea.Value()
			m.textarea.SetHeight(m.textarea.Height() + 1)
			m.textarea.SetValue(currentContent)
		}
	}

	return m, tea.Batch(taCmd, vpCmd)
}

func (m TestCheck) View() string {
	const delimeter = "\n"
	if !m.LastQuestionSuccess {
		return fmt.Sprintf("%s%s%s", m.vpFailed.View(), delimeter, m.textarea.View())
	}
	if len(m.Questions) > 0 {
		return fmt.Sprintf("%s%s%s", m.viewport.View(), delimeter, m.textarea.View())
	} else {
		return fmt.Sprintf("%s%s%s%s%s", m.viewport.View(), delimeter, "Нажми r(reset)/c(choose) для продолжения",
			delimeter, m.textarea.View())
	}
}

func main() {
	program := tea.NewProgram(initializeModel())
	_, err := program.Run()
	if err != nil {
		panic(err)
	}
}
