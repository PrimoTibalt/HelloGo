package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	persistence "primotibalt/checkTests/topicpersistence"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

const (
	paddingToLeft           = 4
	defaultTaPlaceholder    = "Напиши ответ на вопрос"
	addNewQuestionToATopic  = "Добавить новый вопрос в топик"
	testKnowledgeOnTheTopic = "Проверить знания по топику"
	addNewTopic             = "Добавить новый топик"
	removeTopic             = "Удалить топик"
	editTopic               = "Редактировать топик"
	Padding                 = 1 // for some reason it just works and prevent first line of the select from disappearing
)

type mainOption struct {
	name   string
	action func([]persistence.TopicName)
}

var mainOptions []mainOption

type QaPair struct {
	Question string
	Answer   string
}

func init() {
	mainOptions = []mainOption{
		{addNewQuestionToATopic, addNewQuestionToTopicFunc},
		{testKnowledgeOnTheTopic, testKnowledgeOnTheTopicFunc},
		{addNewTopic, addNewTopicFunc},
		{removeTopic, removeTopicFunc},
		{editTopic, editTopicFunc},
	}
}

func main() {
	for {
		var action int
		topics := persistence.RetrieveTopicNames()
		options := make([]huh.Option[int], len(mainOptions))
		for ptr, option := range mainOptions {
			options[ptr] = huh.NewOption(option.name, ptr)
		}
		err := huh.NewForm(huh.NewGroup(
			huh.NewSelect[int]().
				Options(options...).
				Title("Знания - сила. Что делать будем?").
				Value(&action))).
			WithProgramOptions(tea.WithAltScreen()).Run()
		if err != nil {
			if err != huh.ErrUserAborted {
				panic(err)
			}

			os.Exit(0)
		}

		mainOptions[action].action(topics)
	}
}

func testKnowledgeOnTheTopicFunc(topics []persistence.TopicName) {
	selectedTopics := ChooseTopicsForTest(topics)
	if len(selectedTopics) == 0 {
		fmt.Println("Вы не выбрали ничего.")
		return
	}
	RunKnowledgeTest(selectedTopics)
}

func addNewQuestionToTopicFunc(topics []persistence.TopicName) {
	topicToQa := getAllQuestionsForAllTopics(topics)
	selectedTopic, ok := ChooseTopic(topics, topicToQa)
	if !ok {
		return
	}

	RunTopicQuestionAppend(selectedTopic)
}

func addNewTopicFunc(topics []persistence.TopicName) {
	AddNewTopic(topics)
}

func removeTopicFunc(topics []persistence.TopicName) {
	topicToQa := getAllQuestionsForAllTopics(topics)
	RemoveTopic(topics, topicToQa)
}

func editTopicFunc(topics []persistence.TopicName) {
	topicToQA := make(map[persistence.TopicName][]QaPair)
	for _, topic := range topics {
		pairs := []QaPair{}
		for _, qaLine := range persistence.TopicQuestions(topic) {
			qa := strings.Split(qaLine, persistence.DelimeterQuestionAnswer)
			if len(qa) < 2 {
				continue
			}
			pairs = append(pairs, QaPair{qa[0], qa[1]})
		}

		topicToQA[topic] = pairs
	}

	EditTopic(topicToQA)
}

func getAllQuestionsForAllTopics(topics []persistence.TopicName) (topicToQa map[persistence.TopicName][]string) {
	topicToQa = map[persistence.TopicName][]string{}
	for _, topic := range topics {
		qaPairs := persistence.TopicQuestions(topic)
		questions := make([]string, len(qaPairs))
		for questionPtr, qaLine := range qaPairs {
			questions[questionPtr] = strings.Split(qaLine, persistence.DelimeterQuestionAnswer)[0]
		}

		topicToQa[topic] = questions
	}

	return
}

func MoveCursorToTopLeft() {
	command := exec.Command("clear")
	command.Stdout = os.Stdout
	command.Run()
	fmt.Print("\033[H")
}

func GetQuestionsFromTopics(topic persistence.TopicName, topicToQa map[persistence.TopicName][]string) (questions string) {
	var questionsLineSb strings.Builder
	for _, question := range topicToQa[topic] {
		questionsLineSb.WriteString(question + "\n")
	}
	questions = questionsLineSb.String()
	return
}
