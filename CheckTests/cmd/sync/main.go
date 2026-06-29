package main

import (
	persistence "primotibalt/checkTests/topicpersistence"
)

var breaker chan bool

func main() {
	AttachWatcher(persistence.QuestionsDir())
	RunServer()
}
