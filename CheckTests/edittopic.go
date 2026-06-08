package main

import (
	"errors"
	"os"
	persistence "primotibalt/checkTests/topicpersistence"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

type EditTopicModel struct {
	Input    *huh.Text
	Original string
	Select   *huh.Form

	EditingPart EditingPart
	IsEditing   bool

	Content          map[persistence.TopicName]map[string]string
	SelectedTopic    persistence.TopicName
	SelectedQuestion string
}

type EditingPart string

var (
	TopicPart     EditingPart    = "Topic"
	QuestionPart  EditingPart    = "Question"
	AnswerPart    EditingPart    = "Answer"
	selectedValue string         = ""
	border                       = lipgloss.BlockBorder()
	borderStyle   lipgloss.Style = lipgloss.NewStyle().Border(border, true)
)

func (m EditTopicModel) View() string {
	w, h, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		os.Exit(1)
	}
	if m.IsEditing {
		inputTheme := huh.ThemeBase()
		inputTheme.Focused.Base = lipgloss.NewStyle().
			PaddingLeft(Padding).
			PaddingRight(Padding).
			AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Center).
			MaxWidth(w - border.GetLeftSize() - Padding).
			MaxHeight(h - border.GetTopSize() - Padding)
		inputTheme.Focused.Card = inputTheme.Focused.Base
		content := lipgloss.JoinVertical(
			lipgloss.Center,
			renderOriginal(m.Original, w-border.GetLeftSize()-border.GetRightSize()-Padding),
			m.Input.
				WithTheme(inputTheme).
				WithWidth(w-border.GetLeftSize()-border.GetRightSize()-Padding).
				WithHeight(h-border.GetTopSize()-border.GetBottomSize()-4*Padding).
				View())
		return borderStyle.Render(content)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.Select.WithWidth(w-Padding).
			WithHeight(h-Padding).
			View())
}

func (m EditTopicModel) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m EditTopicModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	if !m.IsEditing {
		m.processInputInSelectingMode(msg)
		return m, nil
	}
	if m.IsEditing {
		m.processInputInEditingMode(msg)
		return m, nil
	}
	return m, nil
}

func (m *EditTopicModel) processInputInSelectingMode(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if (msg.String() == "tab" || msg.String() == "enter") && m.EditingPart != AnswerPart {
			m.calculateNextState()
		}
		if msg.String() == "shift+tab" && m.EditingPart != TopicPart {
			m.calculatePrevState()
		}
		if msg.String() == "e" {
			m.updateInputWithNewValue()
		}
		teaModel, _ := m.Select.Update(msg)
		teaSelect, _ := teaModel.(*huh.Form)
		m.Select = teaSelect
	default:
		teaModel, _ := m.Select.Update(msg)
		teaSelect, _ := teaModel.(*huh.Form)
		m.Select = teaSelect
	}
}

func (m *EditTopicModel) processInputInEditingMode(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" || msg.String() == "enter" {
			if msg.String() == "enter" {
				m.applyChanges()
				m.updateSelectWithOptions()
			}
			m.IsEditing = false
			if m.EditingPart == AnswerPart {
				m.calculatePrevState()
			}
			return
		}
	}
	input, _ := m.Input.Update(msg)
	inputText, ok := input.(*huh.Text)
	if ok {
		m.Input = inputText
	}
}

func (m *EditTopicModel) applyChanges() {
	updatedValue := strings.Trim(selectedValue, "\n\t ")
	if updatedValue == "" || strings.Contains(selectedValue, "\n") {
		var errorMsg string
		switch {
		case updatedValue == "":
			errorMsg = "Новое значение не содержит никаких символов, так нельзя"
		case strings.Contains(selectedValue, "\n"):
			errorMsg = "Нельзя переходить на новую строку в тексте нового значения"
		}
		displayErrorNotification(errorMsg)
		return
	}
	switch m.EditingPart {
	case TopicPart:
		if _, ok := m.Content[persistence.TopicName(updatedValue)]; ok {
			return
		}
		initialTopic := persistence.TopicName(m.Original)
		content := m.Content[initialTopic]
		delete(m.Content, initialTopic)
		m.Content[persistence.TopicName(updatedValue)] = content
	case QuestionPart:
		if _, ok := m.Content[persistence.TopicName(updatedValue)]; ok {
			return
		}
		initialQuestion := m.Original
		content := m.Content[m.SelectedTopic][initialQuestion]
		delete(m.Content[m.SelectedTopic], initialQuestion)
		m.Content[m.SelectedTopic][updatedValue] = content
	case AnswerPart:
		m.Content[m.SelectedTopic][m.SelectedQuestion] = updatedValue
	}
}

func (m *EditTopicModel) updateSelectWithOptions() {
	m.Select = huh.NewForm(huh.NewGroup(huh.NewSelect[string]().
		Options(func() []huh.Option[string] {
			options := []huh.Option[string]{}
			switch m.EditingPart {
			case TopicPart:
				for topic := range m.Content {
					options = append(options, huh.NewOption(string(topic), string(topic)))
				}
			case QuestionPart:
				for question := range m.Content[m.SelectedTopic] {
					options = append(options, huh.NewOption(question, question))
				}
			case AnswerPart:
				for _, answer := range m.Content[m.SelectedTopic] {
					options = append(options, huh.NewOption(answer, answer))
				}
			default:
				panic(errors.New("неизвестное состояние для селекта"))
			}
			return options
		}()...).
		Value(&selectedValue).
		Title(getTitleForSelectByState(m.EditingPart))))
	m.Select.NextField()
}

func (m *EditTopicModel) updateInputWithNewValue() {
	m.IsEditing = true
	var note string
	switch m.EditingPart {
	case TopicPart:
		m.SelectedTopic = persistence.TopicName(selectedValue)
		note = selectedValue
	case QuestionPart:
		m.SelectedQuestion = selectedValue
		note = selectedValue
	case AnswerPart:
		note = m.Content[m.SelectedTopic][m.SelectedQuestion]
	}
	m.Original = note
	m.Input = huh.NewText().Value(&selectedValue)
	m.Input.Focus()
}

func (m *EditTopicModel) calculatePrevState() {
	switch m.EditingPart {
	case TopicPart:
	case QuestionPart:
		m.EditingPart = TopicPart
		m.SelectedTopic = ""
	case AnswerPart:
		m.EditingPart = QuestionPart
		m.IsEditing = false
		m.SelectedQuestion = ""
	}
	m.updateSelectWithOptions()
}

func (m *EditTopicModel) calculateNextState() {
	switch m.EditingPart {
	case TopicPart:
		m.EditingPart = QuestionPart
		m.SelectedTopic = persistence.TopicName(selectedValue)
	case QuestionPart:
		m.EditingPart = AnswerPart
		m.SelectedQuestion = selectedValue
	case AnswerPart:
	}
	if m.EditingPart == AnswerPart {
		m.IsEditing = true
		m.SelectedQuestion = selectedValue
		note := m.Content[m.SelectedTopic][m.SelectedQuestion]
		selectedValue = note
		m.Original = note
		m.Input = huh.NewText().Value(&selectedValue)
		m.Input.Focus()
	} else {
		m.updateSelectWithOptions()
	}
}

func EditTopic(topicToQa map[persistence.TopicName][]QaPair) {
	model := initializeEditTopicModel(topicToQa)
	_, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
	if err != nil {
		panic(err)
	}
}

func initializeEditTopicModel(topicToQa map[persistence.TopicName][]QaPair) EditTopicModel {
	m := EditTopicModel{
		EditingPart: TopicPart,
	}
	content := map[persistence.TopicName]map[string]string{}
	for topic, v := range topicToQa {
		content[topic] = map[string]string{}
		for _, qaPair := range v {
			content[topic][qaPair.Question] = qaPair.Answer
		}
	}
	m.Content = content
	m.updateSelectWithOptions()
	return m
}

func getTitleForSelectByState(editingPart EditingPart) string {
	switch editingPart {
	case TopicPart:
		return "Выберите топик для редактирования"
	case QuestionPart:
		return "Выберите вопрос для редактирования"
	case AnswerPart:
		return "Редактируйте ответ на вопрос"
	default:
		panic(errors.New("непредвиденное состояние селекта"))
	}
}

func displayErrorNotification(errorMsg string) {
	theme := huh.ThemeBase()
	theme.Focused.Base = lipgloss.NewStyle().
		Border(border, true).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
	theme.Focused.Card = theme.Focused.Base
	theme.Focused.Card = theme.Focused.Card.Foreground(lipgloss.Color("#ff5555"))
	huh.NewForm(huh.NewGroup(
		huh.NewNote().Title(errorMsg).WithTheme(theme),
	)).
		WithProgramOptions(tea.WithAltScreen()).
		Run()
}

func renderOriginal(original string, maxWidth int) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7571F9")).
		PaddingTop(1).
		Inline(true).
		MaxWidth(maxWidth).
		Render(original)
}
