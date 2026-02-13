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
)

func RetrieveTopicToPathMap() (result map[string]string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	pathToTopics := path.Join(homeDir, PathToQuestionsDir)
	if _, err := os.Stat(pathToTopics); err != nil {
		if err = os.Mkdir(pathToTopics, 0o751); err != nil {
			panic(err)
		}
	}

	questionsDirEntries, getqInfoErr := os.ReadDir(pathToTopics)
	if getqInfoErr != nil {
		panic(getqInfoErr)
	}

	result = make(map[string]string, len(questionsDirEntries))
	for _, dir := range questionsDirEntries {
		result[dir.Name()] = path.Join(pathToTopics, dir.Name())
	}

	return
}

func TopicQuestions(path string) (qaPairs []string) {
	fi, fiErr := os.ReadFile(path)
	if fiErr != nil {
		panic(fiErr)
	}

	content := string(fi)
	return strings.Split(content, "\n")
}
