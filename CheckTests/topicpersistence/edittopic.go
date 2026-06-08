package topicpersistence

import "slices"

func UpdateTopicNames(topics []TopicName) {
	origTopicNames := RetrieveTopicNames()
	topicsToDelete := []TopicName{}
	for _, topic := range origTopicNames {
		if slices.Index(topics, topic) < 0 {
			topicsToDelete = append(topicsToDelete, topic)
		}
	}
}
