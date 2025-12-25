package main

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	retriever "primotibalt/checkTests/questionsRetriever"
)

type Question struct {
	Answer string
	Text   string
}
type TestCheck struct {
	viewport            viewport.Model
	textarea            textarea.Model
	Questions           []Question
	SuccessQuestions    []Question
	FailedQuestions     []Question
	CurrentQuestion     *Question
	LastQuestionSuccess bool
}

func (m TestCheck) Init() tea.Cmd {
	return textarea.Blink
}

func (m TestCheck) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
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
					fmt.Println("You don't have anything to choose")
					return m, tea.Quit
				}
			}

			m.compareInputWithAnswer()
			m.textarea.Reset()

			if m.CurrentQuestion != nil {
				questionIndex := slices.Index(m.Questions, *m.CurrentQuestion)
				m.Questions = slices.Delete(m.Questions, questionIndex, questionIndex+1)
			}

			questionsToAskCount := len(m.Questions)
			if questionsToAskCount > 0 {
				selectedQuestionIndex := rand.Intn(questionsToAskCount)
				m.CurrentQuestion = &m.Questions[selectedQuestionIndex]
				m.viewport.SetContent(m.CurrentQuestion.Text)
			} else {
				m.viewport.SetContent(m.getResultWhenNoMoreQuestions())
				return m, tea.Batch(taCmd, vpCmd)
			}
		}
	}

	return m, tea.Batch(taCmd, vpCmd)
}

func (m TestCheck) View() string {
	if len(m.Questions) > 0 {
		return fmt.Sprintf("%s%s%s", m.textarea.View(), "\n", m.viewport.View())
	} else {
		return fmt.Sprintf("%s%s%s", m.viewport.View(), "\n", "Нажми r(reset)/c(choose) для продолжения")
	}
}

func main() {
	program := tea.NewProgram(initializeModel())
	_, err := program.Run()
	if err != nil {
		panic(err)
	}
}

func (m *TestCheck) getResultWhenNoMoreQuestions() (result string) {
	sb := strings.Builder{}
	if len(m.FailedQuestions) > 0 {
		sb.WriteString(
			fmt.Sprintf("Неправильно ответил на %d вопросов из %d.\n",
				len(m.FailedQuestions),
				len(m.FailedQuestions)+len(m.SuccessQuestions)))
		sb.WriteString("Заваленные вопросы:\n")
		for _, question := range m.FailedQuestions {
			sb.WriteString(fmt.Sprintf("%+v\n", question))
		}
	} else {
		sb.WriteString("Вы не завалили ни одного вопроса. Молодец!\n")
	}

	result = sb.String()
	return
}

func (m *TestCheck) compareInputWithAnswer() {
	if strings.Trim(m.textarea.Value(), " \n\r") == m.CurrentQuestion.Answer {
		m.SuccessQuestions = append(m.SuccessQuestions, *m.CurrentQuestion)
		m.LastQuestionSuccess = true
	} else {
		m.FailedQuestions = append(m.FailedQuestions, *m.CurrentQuestion)
		m.LastQuestionSuccess = false
	}
}

func initializeModel() (testCheckModel TestCheck) {
	topics := retriever.RetrieveTopicToPathMap()
	fmt.Println("Выбери топик для теста:")
	orderedMapOfTopics := make(map[int]string, len(topics))
	i := 1
	for topicName := range topics {
		orderedMapOfTopics[i] = topicName
		fmt.Printf("%d. %s\n", i, topicName)
		i++
	}

	var num int
	fmt.Scan(&num)
	qaPairs := retriever.TopicQuestions(topics[orderedMapOfTopics[num]])
	questions := []Question{}
	for _, pair := range qaPairs {
		if !strings.Contains(pair, "//") {
			continue
		}
		fmt.Printf("%s\n", pair)
		qa := strings.Split(pair, "//")
		answer := qa[1]
		question := qa[0]
		questions = append(questions, Question{answer, question})
	}
	currentQuestion := &questions[rand.Intn(len(questions))]

	textAreaModel := textarea.New()
	textAreaModel.Focus()
	textAreaModel.Placeholder = "Напиши ответ на вопрос"
	textAreaModel.KeyMap = textarea.DefaultKeyMap
	textAreaModel.FocusedStyle.CursorLine = lipgloss.NewStyle()

	vpModel := viewport.New(100, 8)
	vpModel.SetContent(currentQuestion.Text)
	testCheckModel = TestCheck{
		vpModel,
		textAreaModel,
		questions,
		[]Question{},
		[]Question{},
		currentQuestion,
		false,
	}
	return
}
