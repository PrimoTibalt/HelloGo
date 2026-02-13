package main

import (
	"os"
	persistence "primotibalt/checkTests/topicPersistence"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

type AppendQuestion struct {
	Questions   []Question
	LeftPanelVp viewport.Model
	QaForm      *huh.Form
	QuestionSet bool
	PathToFile  string
}

var (
	questionInputDesc string = "Введите вопрос"
	answerInputDesc   string = "Введите ответ на вопрос"
)

func ChooseTopic(topics map[string]string) (selectedPath string) {
	options := make([]huh.Option[string], len(topics))
	var i int
	for topic, path := range topics {
		options[i] = huh.NewOption(topic, path)
		i++
	}
	err := huh.NewForm(huh.NewGroup(huh.NewSelect[string]().
		Title("Выбери топик для добавления вопросов:").
		Options(options...).
		Value(&selectedPath))).
		Run()
	if err != nil {
		if err == huh.ErrUserAborted {
			os.Exit(0)
		}
	}

	return
}

func RunTopicQuestionAppend(selectedPath string) {
	questions := []Question{}
	qaPairs := persistence.TopicQuestions(selectedPath)
	for _, pair := range qaPairs {
		if !strings.Contains(pair, persistence.DelimeterQuestionAnswer) {
			continue
		}

		qa := strings.Split(pair, persistence.DelimeterQuestionAnswer)
		question := qa[0]
		answer := qa[1]
		questions = append(questions, Question{answer, question})
	}

	model := initializeQuestionAppendModel(questions)
	model.PathToFile = selectedPath
	program := tea.NewProgram(model)
	_, err := program.Run()
	if err != nil {
		panic(err)
	}
}

func (appendQuestionModel AppendQuestion) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Left, appendQuestionModel.LeftPanelVp.View(), appendQuestionModel.QaForm.View())
}

func (appendQuestionModel AppendQuestion) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return appendQuestionModel, tea.Batch(tea.ClearScreen, tea.Quit)
		case tea.KeyEnter:
			if appendQuestionModel.QuestionSet {
				answer := appendQuestionModel.QaForm.GetFocusedField().GetValue()
				appendQuestionModel.QaForm.PrevField()
				question := appendQuestionModel.QaForm.GetFocusedField().GetValue()
				appendQuestionModel.Questions = append(appendQuestionModel.Questions, Question{Text: question.(string), Answer: answer.(string)})
				appendQuestionModel.QuestionSet = false
				appendQuestionModel.QaForm = getQaForm(appendQuestionModel.QaForm.Help().Width)
				persistence.AppendQuestionsToTopic(appendQuestionModel.PathToFile, map[string]string{question.(string): answer.(string)})
			} else {
				appendQuestionModel.QaForm.NextField()
				appendQuestionModel.QuestionSet = true
			}
		case tea.KeyCtrlZ:
			if appendQuestionModel.QuestionSet {
				appendQuestionModel.QaForm.PrevField()
				appendQuestionModel.QuestionSet = false
			}
		}
	}
	sb := strings.Builder{}
	for _, item := range appendQuestionModel.Questions {
		sb.WriteString("Q:" + item.PrintText(appendQuestionModel.LeftPanelVp.Width-2) + "\n")
		sb.WriteString("A:" + item.PrintAnswer(appendQuestionModel.LeftPanelVp.Width-2) + "\n\n")
	}
	appendQuestionModel.LeftPanelVp.SetContent(sb.String())
	appendQuestionModel.QaForm.Update(msg)
	appendQuestionModel.LeftPanelVp.Update(msg)

	return appendQuestionModel, nil
}

func (appendQuestionModel AppendQuestion) Init() tea.Cmd {
	return tea.ClearScreen
}

func initializeQuestionAppendModel(questions []Question) (appendQuestionModel AppendQuestion) {
	width, height, termSizeErr := term.GetSize(os.Stdout.Fd())
	rightPanelWidth := width * 3 / 7
	leftPanelWidth := width - rightPanelWidth
	if termSizeErr != nil {
		width = 80
	}

	appendQuestionModel = AppendQuestion{
		LeftPanelVp: viewport.New(leftPanelWidth, height),
		Questions:   questions,
		QaForm:      getQaForm(rightPanelWidth),
	}

	return
}

func getQaForm(rightPanelWidth int) *huh.Form {
	questionInput := huh.NewText().Description(questionInputDesc)
	questionInput.Focus()
	answerInput := huh.NewText().Description(answerInputDesc)
	return huh.NewForm(huh.NewGroup(
		questionInput,
		answerInput,
	)).WithWidth(rightPanelWidth)
}
