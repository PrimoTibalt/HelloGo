package main

import (
	persistence "primotibalt/checkTests/topicpersistence"
)

func main() {
	var (
		breaker  = make(chan bool)
		notifier = make(chan bool)
	)
	AttachWatcher(persistence.QuestionsDir(), breaker, notifier)
	RunServer(breaker, notifier)
}
