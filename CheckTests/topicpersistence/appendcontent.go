package topicpersistence

import (
	"os"
	"path"
)

func AppendQuestionsToTopic(topicName TopicName, questionsToAdd map[string]string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	filepath := path.Join(homeDir, PathToQuestionsDir, string(topicName))
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	if err = appendQuestionToFile(file, questionsToAdd); err != nil {
		panic(err)
	}
}

func appendQuestionToFile(file *os.File, questionsToAdd map[string]string) error {
	for q, a := range questionsToAdd {
		_, err := file.Write(
			[]byte(q +
				DelimeterQuestionAnswer +
				a +
				"\n"))
		if err != nil {
			return err
		}
	}

	return nil
}
