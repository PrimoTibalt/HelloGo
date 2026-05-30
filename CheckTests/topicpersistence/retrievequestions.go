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

func RetrieveTopicToPathMap() (topics []string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	pathToTopics := path.Join(homeDir, PathToQuestionsDir)
	if _, err := os.Stat(pathToTopics); err != nil {
		if err = os.MkdirAll(pathToTopics, directoryPermissions); err != nil {
			panic(err)
		}
	}

	questionsDirEntries, getqInfoErr := os.ReadDir(pathToTopics)
	if getqInfoErr != nil {
		panic(getqInfoErr)
	}

	topics = make([]string, len(questionsDirEntries))
	for ptr, dir := range questionsDirEntries {
		topics[ptr] = dir.Name()
	}

	return
}

func TopicQuestions(topic string) (qaPairs []string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	pathToTopic := path.Join(homeDir, PathToQuestionsDir, topic)
	fi, fiErr := os.ReadFile(pathToTopic)
	if fiErr != nil {
		panic(fiErr)
	}

	content := string(fi)
	return strings.Split(content, "\n")
}
