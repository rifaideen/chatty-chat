package config

import (
	"pkg/utils"
	"strings"
)

type Config struct {
	Brokers       []string // kafka brokers
	Group         string   // kafka group
	TopicConsumer string   // kafka consumer topic
	TopicProducer string   // kafka producer topic
	Dsn           string   // database dsn
}

func Load() *Config {
	brokers := utils.GetEnv("KAFKA_BROKERS", "localhost:9092")

	return &Config{
		Brokers:       strings.Split(brokers, ","),
		Group:         utils.GetEnv("KAFKA_GROUP", "persistence-group"),
		TopicConsumer: utils.GetEnv("KAFKA_TOPIC_CONSUMER", "chat"),         // consume chat messages
		TopicProducer: utils.GetEnv("KAFKA_TOPIC_PRODUCER", "notification"), // produce notification messages
		Dsn:           utils.GetEnv("DSN", ""),
	}
}
