package main

import (
	"os"

	retriever "primotibalt/checkTests/questionsRetriever"

	"github.com/charmbracelet/huh"
)

const (
	paddingToLeft           = 4
	defaultTaPlaceholder    = "Напиши ответ на вопрос"
	delimeterQuestionAnswer = "/!/"
	addNewQuestionToATopic  = "Добавить новый вопрос в топик"
	testKnowledgeOnTheTopic = "Проверить знания по топику"
)

var mainOptions map[string]func(map[string]string)

func main() {
	for {
		var action string
		topics := retriever.RetrieveTopicToPathMap()
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

func addNewQuestionToATopicFunc(topics map[string]string) {
}

func init() {
	mainOptions = map[string]func(map[string]string){
		addNewQuestionToATopic:  addNewQuestionToATopicFunc,
		testKnowledgeOnTheTopic: testKnowledgeOnTheTopicFunc,
	}
}
