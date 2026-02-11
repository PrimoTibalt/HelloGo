package main

import (
	"os"

	persistence "primotibalt/checkTests/topicPersistence"

	"github.com/charmbracelet/huh"
)

const (
	paddingToLeft           = 4
	defaultTaPlaceholder    = "Напиши ответ на вопрос"
	addNewQuestionToATopic  = "Добавить новый вопрос в топик"
	testKnowledgeOnTheTopic = "Проверить знания по топику"
)

var mainOptions map[string]func(map[string]string)

func main() {
	for {
		var action string
		topics := persistence.RetrieveTopicToPathMap()
		options := make([]huh.Option[string], len(mainOptions))
		idx := 0
		for key := range mainOptions {
			options[idx] = huh.NewOption(key, key)
			idx++
		}
		err := huh.NewForm(huh.NewGroup(huh.NewSelect[string]().
			Options(options...).
			Value(&action))).Run()
		if err != nil {
			if err != huh.ErrUserAborted {
				panic(err)
			}

			os.Exit(0)
		}

		mainOptions[action](topics)
	}
}

func testKnowledgeOnTheTopicFunc(topics map[string]string) {
	selectedPaths := ChooseTopicForTest(topics)
	RunKnowledgeTest(selectedPaths)
}

func addNewQuestionToTopicFunc(topics map[string]string) {
	selectedPath := ChooseTopic(topics)
	RunTopicQuestionAppend(selectedPath)
}

func init() {
	mainOptions = map[string]func(map[string]string){
		addNewQuestionToATopic:  addNewQuestionToTopicFunc,
		testKnowledgeOnTheTopic: testKnowledgeOnTheTopicFunc,
	}
}
