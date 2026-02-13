package topicpersistence

import "os"

func AppendQuestionsToTopic(pathToFile string, questionsToAdd map[string]string) {
	file, err := os.OpenFile(pathToFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	for q, a := range questionsToAdd {
		_, err = file.Write(
			[]byte(q +
				DelimeterQuestionAnswer +
				a +
				"\n"))
		if err != nil {
			panic(err)
		}
	}
}
