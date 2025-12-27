// Package questionsretriever
package questionsretriever

import (
	"os"
	"path"
	"strings"
)

func RetrieveTopicToPathMap() (result map[string]string) {
	curDir, getCurDirErr := os.Executable()
	pathToTopics := path.Join(curDir, "../Questions")
	if getCurDirErr != nil {
		panic(getCurDirErr)
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
