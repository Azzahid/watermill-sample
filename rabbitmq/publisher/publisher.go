package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
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
	pub, err := amqp.NewPublisher(amqpConfig, watermill.NewStdLogger(false, false))
	if err != nil {
		log.Fatalf("Failed to create publisher: %v", err)
	}
	defer pub.Close()

	for i := 0; true; i++ {
		tracker := map[string]string{
			"id":           strconv.Itoa(i),
			"name":         "Jane Doe",
			"phone_number": "081234567890",
			"status":       "active",
		}

		jsonMessage, err := json.Marshal(tracker)
		if err != nil {
			log.Fatalf("Failed to encode tracker: %v", err)
		}

		// Create a new message
		msg := message.NewMessage(
			watermill.NewUUID(),
			jsonMessage,
		)

		topic := "sample-user"
		if err := pub.Publish(topic, msg); err != nil {
			log.Fatalf("Failed to publish message: %v", err)
		}

		log.Printf("Successfully published message to topic: %s", topic)
		time.Sleep(1 * time.Second)
	}
}
