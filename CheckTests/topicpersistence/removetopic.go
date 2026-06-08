// Package topicpersistence
package topicpersistence

import (
	"os"
	"path"
)

func RemoveTopic(topicName TopicName) (ok bool, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	filepath := path.Join(homeDir, PathToQuestionsDir, string(topicName))
	err = os.Remove(filepath)
	if err != nil {
		return false, err
	}

	return true, nil
}
