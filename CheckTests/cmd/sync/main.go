package main

import (
	persistence "primotibalt/checkTests/topicpersistence"
)

func main() {
	AttachWatcher(persistence.QuestionsDir())
	RunServer()
}
