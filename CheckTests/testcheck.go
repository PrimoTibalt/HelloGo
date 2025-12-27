package main

import (
	"fmt"
	"math/rand"
	"os"
	retriever "primotibalt/checkTests/questionsRetriever"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

type Question struct {
	Answer string
	Text   string
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

func (m *TestCheck) isInputAndAnswerEqual() bool {
	return strings.Trim(m.textarea.Value(), " \n\r") == m.CurrentQuestion.Answer
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
		if !strings.Contains(pair, delimeterQuestionAnswer) {
			continue
		}

		qa := strings.Split(pair, delimeterQuestionAnswer)
		question := qa[0]
		answer := qa[1]
		questions = append(questions, Question{answer, question})
	}

	width, _, termSizeErr := term.GetSize(os.Stdout.Fd())
	if termSizeErr != nil {
		width = 80
	}

	currentQuestion := &questions[rand.Intn(len(questions))]

	taModel := textarea.New()
	taModel.Focus()
	taModel.SetHeight(1)
	taModel.KeyMap = textarea.DefaultKeyMap
	taModel.Placeholder = defaultTaPlaceholder
	taModel.FocusedStyle.CursorLine = lipgloss.NewStyle()
	taModel.ShowLineNumbers = false
	taModel.SetWidth(width - paddingToLeft)

	vpModel := viewport.New(width-paddingToLeft, 3)
	vpModel.SetContent(currentQuestion.Text)

	vpFailedModel := viewport.New(width-paddingToLeft, 7)
	vpFailedModelBorder := lipgloss.NormalBorder()
	vpFailedModel.Style = vpFailedModel.Style.Border(vpFailedModelBorder).
		BorderForeground(lipgloss.Color("9"))

	testCheckModel = TestCheck{
		vpModel,
		taModel,
		vpFailedModel,
		questions,
		[]Question{},
		[]Question{},
		currentQuestion,
		true,
	}
	return
}
