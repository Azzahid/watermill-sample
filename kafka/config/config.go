package config

import (
	"os"
	"strings"
)

type KafkaConfig struct {
	Brokers     []string
	ClientID    string
	Credentials Credentials
}

func GetKafkaConfig() KafkaConfig {
	// Get brokers from env (comma-separated list)
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "kafka:29092"
	}

	return KafkaConfig{
		Brokers:  strings.Split(brokers, ","),
		ClientID: "test",
		Credentials: Credentials{
			Protocol:  "PLAINTEXT",
			Mechanism: "PLAIN",
		},
	}
}
