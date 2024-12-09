package main

import (
	"context"
	"encoding/json"
	"kafka/config"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewSubscriber() (message.Subscriber, error) {
	kafkaConfig := config.GetKafkaConfig()
	saramaConfig := kafka.DefaultSaramaSyncPublisherConfig()

	err := config.AddCredentials(saramaConfig, kafkaConfig.Credentials)
	if err != nil {
		return nil, err
	}

	subscriberConfig := kafka.SubscriberConfig{
		Brokers:               kafkaConfig.Brokers,
		Unmarshaler:           kafka.DefaultMarshaler{},
		OverwriteSaramaConfig: saramaConfig,
		ConsumerGroup:         "group-1",
	}

	subscriber, err := kafka.NewSubscriber(subscriberConfig, watermill.NewStdLogger(false, true))
	if err != nil {
		return nil, err
	}

	return subscriber, nil
}

func main() {
	subscriber, err := NewSubscriber()
	if err != nil {
		panic(err)
	}

	messages, err := subscriber.Subscribe(context.Background(), "test")
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()
	for message := range messages {
		var jsonMessage map[string]interface{}
		err = json.Unmarshal(message.Payload, &jsonMessage)
		if err != nil {
			log.Printf("error unmarshalling message: %v\n", err)
			continue
		}
		log.Printf("received message: %v\n", jsonMessage)
		message.Ack()
	}
}
