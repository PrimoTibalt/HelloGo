// Package topicpersistence
package topicpersistence

import (
	"os"
	"path"
)

func AddNewTopic(topicName string) (filepath string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	filepath = path.Join(homeDir, PathToQuestionsDir, topicName)
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
	return
}
