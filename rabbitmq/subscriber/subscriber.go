package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
)

func main() {
	fmt.Println("Starting RabbitMQ Publisher")
	// Initialize Kafka configuration
	RabbitURI := os.Getenv("RABBITMQ_URI")
	if RabbitURI == "" {
		RabbitURI = "amqp://guest:guest@rabbitmq-sample:5672/"
	}
	amqpConfig := amqp.NewDurablePubSubConfig(RabbitURI, amqp.GenerateQueueNameTopicName)
	// Create a new publisher
	subs, err := amqp.NewSubscriber(amqpConfig, watermill.NewStdLogger(false, false))
	if err != nil {
		log.Fatalf("Failed to create subscriber: %v", err)
	}
	defer subs.Close()

	sub, err := subs.Subscribe(context.Background(), "sample-user")
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	for msg := range sub {
		var jsonMessage map[string]interface{}
		err = json.Unmarshal(msg.Payload, &jsonMessage)
		if err != nil {
			log.Printf("error unmarshalling message: %v\n", err)
			continue
		}
		log.Printf("received message: %v\n", jsonMessage)
		msg.Ack()
	}
}
