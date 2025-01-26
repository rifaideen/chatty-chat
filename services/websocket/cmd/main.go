package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pkg/kafka"
	"websocket/internal/config"
	"websocket/internal/service"

	"github.com/IBM/sarama"
)

var PORT = ":8003"

func main() {
	service := bootstrap()

	// create cancelable context and defer cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start listening to handle the websocket connections regisrations and messages
	go service.Listen(ctx)

	go func() {
		err := service.Consume()

		if err != nil {
			log.Fatal(err)
		}
	}()

	http.HandleFunc("GET /ws", service.ServeWS)

	fmt.Printf("websocket listening on http://localhost%s\n", PORT)

	log.Fatal(http.ListenAndServe(PORT, nil))
}

// bootstrap the application and return service instance
func bootstrap() service.WebsocketService {
	config := config.Load()

	cfg := sarama.NewConfig()

	// producer and consumer configurations
	producer(cfg)
	consumer(cfg)

	// configure kafka producer and consumer
	producer, consumer := setupKafka(config.Brokers, config.Group, cfg)

	service := service.New(
		consumer,
		producer,
		[]string{
			config.ProducerTopic,
			config.ConsumerTopic,
		},
		config.AuthServiceUrl,
		config.OllamaServiceUrl,
	)

	return service
}

// kafka consumer configurations
func consumer(config *sarama.Config) {
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
}

// kafka producer configurations
func producer(config *sarama.Config) {
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
}

func setupKafka(brokers []string, group string, config *sarama.Config) (kafka.Producer, kafka.Consumer) {
	client := kafka.New(brokers, group, config)

	producer, err := client.SetupProducer()

	if err != nil {
		log.Fatal(err)
	}

	consumer, err := client.SetupConsumer()

	if err != nil {
		log.Fatal(err)
	}

	return producer, consumer
}
