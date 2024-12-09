package main

import (
	"encoding/json"
	"kafka/config"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

func NewPublisher() (message.Publisher, error) {
	kafkaConfig := config.GetKafkaConfig()
	saramaConfig := kafka.DefaultSaramaSyncPublisherConfig()

	err := config.AddCredentials(saramaConfig, kafkaConfig.Credentials)
	if err != nil {
		return nil, err
	}

	saramaConfig.ClientID = kafkaConfig.ClientID

	publisherConfig := kafka.PublisherConfig{
		Brokers:               kafkaConfig.Brokers,
		Marshaler:             kafka.DefaultMarshaler{},
		OverwriteSaramaConfig: saramaConfig,
	}

	// Create a new publisher with the config
	publisher, err := kafka.NewPublisher(publisherConfig, watermill.NewStdLogger(false, true))
	if err != nil {
		return nil, err
	}

	return publisher, nil
}

func main() {
	publisher, err := NewPublisher()
	if err != nil {
		panic(err)
	}
	defer publisher.Close()
	for i := 0; true; i++ {
		json, err := json.Marshal(map[string]interface{}{
			"category": "test",
			"value":    i,
		})
		if err != nil {
			log.Printf("error marshalling message: %v\n", err)
			continue
		}

		err = publisher.Publish("test", message.NewMessage(uuid.New().String(), json))
		if err != nil {
			log.Printf("error publishing message: %v\n", err)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
