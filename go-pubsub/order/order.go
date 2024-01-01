package order

import (
	"cloud.google.com/go/pubsub"
)

func publish(client *pubsub.Client, topic *pubsub.Topic, message string) error {
	result := topic.Publish(client.Context(), &pubsub.Message{
		Data: []byte(message),
	})
	_, err := result.Get(client.Context())
	return err
}
