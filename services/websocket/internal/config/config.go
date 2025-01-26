package config

import (
	"pkg/utils"
	"strings"
)

type Config struct {
	Brokers          []string // kafka brokers
	Group            string   // kafka group
	ConsumerTopic    string   // topic to consume from kafka
	ProducerTopic    string   // topic to produce into kafka
	AuthServiceUrl   string   // auth service url
	OllamaServiceUrl string   // ollama service url
}

func Load() *Config {
	brokers := utils.GetEnv("KAFKA_BROKERS", "localhost:9092")

	return &Config{
		Brokers:          strings.Split(brokers, ","),
		Group:            utils.GetEnv("KAFKA_GROUP", "websocket-group"),
		ConsumerTopic:    utils.GetEnv("KAFKA_TOPIC_CONSUMER", "notification"), // consumes notification topic
		ProducerTopic:    utils.GetEnv("KAFKA_TOPIC_PRODUCER", "chat"),         // produces chat topic
		AuthServiceUrl:   utils.GetEnv("AUTH_SERVICE_URL", "http://auth-service:8001"),
		OllamaServiceUrl: utils.GetEnv("OLLAMA_SERVICE_URL", "http://ollama_service:11434"),
	}
}
