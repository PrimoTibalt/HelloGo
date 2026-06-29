// Package topicpersistence
package topicpersistence

import (
	"os"
	"path"
	"strings"
)

const (
	DelimeterQuestionAnswer = "/!/"
	PathToQuestionsDir      = "/.local/share/primotibalt/Questions"
	directoryPermissions    = 0o755
)

type TopicName string

func RetrieveTopicNames() (topics []TopicName) {
	pathToTopics := QuestionsDir()
	if _, err := os.Stat(pathToTopics); err != nil {
		if err = os.MkdirAll(pathToTopics, directoryPermissions); err != nil {
			panic(err)
		}
	}

	questionsDirEntries, getqInfoErr := os.ReadDir(pathToTopics)
	if getqInfoErr != nil {
		panic(getqInfoErr)
	}

	topics = make([]TopicName, len(questionsDirEntries))
	for ptr, dir := range questionsDirEntries {
		topics[ptr] = TopicName(dir.Name())
	}

	return
}

func TopicQuestions(topic TopicName) (qaPairs []string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	pathToTopic := path.Join(homeDir, PathToQuestionsDir, string(topic))
	fi, fiErr := os.ReadFile(pathToTopic)
	if fiErr != nil {
		panic(fiErr)
	}

	content := string(fi)
	return strings.Split(content, "\n")
}

func QuestionsDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return path.Join(homeDir, PathToQuestionsDir)
}
