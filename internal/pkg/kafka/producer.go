package kafka

import "github.com/IBM/sarama"

// embeds sarama.AsyncProducer to mock it later
type Producer interface {
	sarama.AsyncProducer
}
