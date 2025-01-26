package kafka

import "github.com/IBM/sarama"

type Service interface {
	SetupConsumer() (Consumer, error)
	SetupProducer() (Producer, error)
}

type KafkaService struct {
	brokers []string
	group   string
	config  *sarama.Config
}

func New(brokers []string, group string, config *sarama.Config) Service {
	return &KafkaService{
		brokers: brokers,
		group:   group,
		config:  config,
	}
}

func (s *KafkaService) SetupConsumer() (Consumer, error) {
	return sarama.NewConsumerGroup(s.brokers, s.group, s.config)
}

func (s *KafkaService) SetupProducer() (Producer, error) {
	return sarama.NewAsyncProducer(s.brokers, s.config)
}
