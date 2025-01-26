package kafka

import "github.com/IBM/sarama"

// embeds sarama.ConsumerGroup to mock it later
type Consumer interface {
	sarama.ConsumerGroup
}
