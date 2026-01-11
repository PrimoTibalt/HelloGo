package main

import (
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

type Question struct {
	Answer string
	Text   string
}

func (q *Question) PrintAnswer(consoleWidth int) string {
	return print(q.Answer, consoleWidth)
}

func (q *Question) PrintText(consoleWidth int) string {
	return print(q.Text, consoleWidth)
}

func print(text string, consoleWidth int) (splitText string) {
	var builder strings.Builder
	var counter int
	var index int
	runes := []rune(text)
	for counter < len(runes) {
		if index == consoleWidth {
			index = index - consoleWidth
			builder.WriteRune('\n')
		}
		builder.WriteRune(runes[counter])
		index++
		counter++
	}

	splitText = builder.String()
	return
}

type TestCheck struct {
	viewport            viewport.Model
	textarea            textarea.Model
	vpFailed            viewport.Model
	Questions           []Question
	SuccessQuestions    []Question
	FailedQuestions     []Question
	CurrentQuestion     *Question
	LastQuestionSuccess bool
}

const (
	smallAnswerLen           = 10
	mediumAnswerLen          = 15
	bigAnswerLen             = 25
	paragraphAnswerLen       = 100
	poemAnswerLength         = 250
	dissertationAnswerLength = 1000
)

func (m *TestCheck) setContentWhenNoMoreQuestions() {
	sb := strings.Builder{}
	if len(m.FailedQuestions) > 0 {
		sb.WriteString(
			fmt.Sprintf("Неправильно ответил на %d вопросов из %d.\n",
				len(m.FailedQuestions),
				len(m.FailedQuestions)+len(m.SuccessQuestions)))
		sb.WriteString("Заваленные вопросы:\n")
		for _, question := range m.FailedQuestions {
			sb.WriteString(fmt.Sprintf("Вопрос: %s\nОтвет: %s\n", question.Text, question.Answer))
		}
	} else {
		sb.WriteString("Вы не завалили ни одного вопроса. Молодец!\n")
	}

	resultString := sb.String()
	m.viewport.SetContent(resultString)
	m.viewport.Height = strings.Count(resultString, "\n") + 1
}

func (m *TestCheck) computeDistanceBetweenInputAndAnswer() (distance int) {
	distance = Ld(strings.Trim(m.textarea.Value(), " \n\r"), m.CurrentQuestion.Answer)
	return
}

func (m *TestCheck) IsInputAndAnswerEqual(distance int) (result bool) {
	answerLen := utf8.RuneCountInString(m.CurrentQuestion.Answer)
	if answerLen <= smallAnswerLen {
		result = distance <= 1
	} else if answerLen <= mediumAnswerLen {
		result = distance <= 2
	} else if answerLen <= bigAnswerLen {
		result = distance <= 4
	} else if answerLen <= paragraphAnswerLen {
		result = distance <= 8
	} else if answerLen <= poemAnswerLength {
		result = distance <= 16
	} else if answerLen <= dissertationAnswerLength {
		result = distance <= 32
	} else {
		result = distance < 100
	}

	return
}

func initializeModel(questions []Question) (testCheckModel TestCheck) {
	width, height, termSizeErr := term.GetSize(os.Stdout.Fd())
	rightPanelWidth := width * 3 / 7
	leftPanelWidth := width - rightPanelWidth
	if termSizeErr != nil {
		width = 80
	}

	currentQuestion := questions[rand.Intn(len(questions))]

	taModel := textarea.New()
	taModel.Focus()
	taModel.SetHeight(2)
	taModel.KeyMap = textarea.DefaultKeyMap
	taModel.Placeholder = defaultTaPlaceholder
	taModel.FocusedStyle.CursorLine = lipgloss.NewStyle()
	taModel.ShowLineNumbers = false
	taModel.SetWidth(leftPanelWidth)

	vpModel := viewport.New(leftPanelWidth, 8)
	vpModel.KeyMap = viewport.KeyMap{}
	vpModel.KeyMap.Down = key.NewBinding(key.WithKeys("down"))
	vpModel.KeyMap.Up = key.NewBinding(key.WithKeys("up"))
	vpModel.Style = vpModel.Style.Border(
		lipgloss.NormalBorder(),
		true, true).
		BorderForeground(lipgloss.Color("6")).
		Foreground(lipgloss.Color("180"))

	vpModel.SetContent(currentQuestion.Text)

	vpFailedModel := viewport.New(rightPanelWidth, height-2)
	vpFailedModel.KeyMap = viewport.KeyMap{}
	vpFailedModel.KeyMap.Down = key.NewBinding(key.WithKeys("down"))
	vpFailedModel.KeyMap.Up = key.NewBinding(key.WithKeys("up"))
	vpFailedModel.Style = vpFailedModel.Style.Border(
		lipgloss.NormalBorder(),
		true, false, true, false).
		BorderForeground(lipgloss.Color("52"))

	testCheckModel = TestCheck{
		vpModel,
		taModel,
		vpFailedModel,
		questions,
		[]Question{},
		[]Question{},
		&currentQuestion,
		true,
	}
	return
}

func (m *TestCheck) prepareNextQuestion() {
	questionIndex := slices.Index(m.Questions, *m.CurrentQuestion)
	m.Questions = slices.Delete(m.Questions, questionIndex, questionIndex+1)

	questionsToAskCount := len(m.Questions)
	if questionsToAskCount > 0 {
		selectedQuestionIndex := rand.Intn(questionsToAskCount)
		m.CurrentQuestion = &m.Questions[selectedQuestionIndex]
		m.viewport.SetContent(m.CurrentQuestion.PrintText(m.viewport.Width))
	} else {
		m.setContentWhenNoMoreQuestions()
	}
}

func (m *TestCheck) prepareFailedVpContent(yourAnswer string) {
	var builder strings.Builder
	for l := range strings.SplitSeq(m.vpFailed.View(), "\n") {
		if l != "" {
			builder.WriteString(l + "\n")
		}
	}

	m.vpFailed.SetContent(
		print("V: "+m.CurrentQuestion.Answer, m.vpFailed.Width-1) + "\n" +
			print("X: "+yourAnswer, m.vpFailed.Width-1) + "\n" +
			builder.String())
}
