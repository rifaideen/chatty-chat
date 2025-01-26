package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"persistence/internal/config"
	"persistence/internal/service"
	"pkg/kafka"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"

	_ "github.com/lib/pq"
)

func main() {
	// load configuration
	config := config.Load()

	// configure logger
	logger()

	// configure database
	db := database(config.Dsn)
	defer db.Close()

	// migrate database
	migrate(db)

	// create kafka configuration
	cfg := sarama.NewConfig()

	// configure producer and consumer
	configure(cfg)

	// setup kafka consumer
	producer, consumer := setup(config.Brokers, config.Group, cfg)

	// close consumer and producer
	defer consumer.Close()
	defer producer.Close()

	// create cancelable context and cancel when interrupt signal is received
	ctx, cancel := context.WithCancel(context.Background())

	// create wait group to wait for goroutines to finish
	var wg sync.WaitGroup
	wg.Add(2)

	// create channel to receive interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// create persistence service
	service := service.New(producer, consumer, db)

	// wait for interrupt signal and cancel context
	go func() {
		<-quit
		cancel()
	}()

	// start consuming messages
	go func() {
		defer wg.Done()
		service.Consume(ctx, config.TopicConsumer)
	}()

	// start producing messages every minute
	go func() {
		defer wg.Done()

		// start the time ticker and produce messages every minute
		t := time.NewTicker(time.Minute)
		msg := &sarama.ProducerMessage{
			Topic: config.TopicProducer,
		}

		for {
			select {
			case <-t.C:
				msg.Value = sarama.StringEncoder("Notification: " + time.Now().Format("2006-01-02 03:04 PM"))

				slog.Info("sending notification message", "topic", config.TopicProducer, "message", msg)

				service.Produce(ctx, msg)
			case <-ctx.Done():
				return
			}
		}
	}()

	slog.Info("persistence service started listening on topic.", "consume", config.TopicConsumer, "produce", config.TopicProducer)
	// wait for goroutines to finish
	wg.Wait()
}

func logger() {
	// configure logger with JSON formatter and set as default
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	slog.SetDefault(log)
}

// kafka consumer configurations
func configure(config *sarama.Config) {
	// consumer configurations
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// producer configurations
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
}

func setup(brokers []string, group string, config *sarama.Config) (kafka.Producer, kafka.Consumer) {
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

func database(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		slog.Warn("error connecting to database", "error", err)
		os.Exit(1)
	}

	return db
}

func migrate(db *sql.DB) {
	// create table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS
			messages (
				id BIGSERIAL PRIMARY KEY,
				message TEXT NOT NULL,
				topic VARCHAR(255) NOT NULL,
				partition INT NOT NULL,
				ofset BIGINT NOT NULL,
				timestamp TIMESTAMP NOT NULL,
				created_at TIMESTAMP NOT NULL DEFAULT NOW ()
			);

		CREATE INDEX IF NOT EXISTS idx_messages_topic_timestamp ON messages (topic, timestamp);
	`)

	if err != nil {
		slog.Warn("error creating table", "error", err)
		os.Exit(1)
	} else {
		slog.Info("database migration completed")
	}
}
