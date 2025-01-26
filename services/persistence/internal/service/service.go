package service

import (
	"context"
	"log/slog"
	"persistence/internal/model"
	"pkg/db"
	"pkg/kafka"
	"sync"

	"github.com/IBM/sarama"
)

type Service interface {
	Consume(ctx context.Context, topic string) error
	Produce(ctx context.Context, msg *sarama.ProducerMessage) error
}

type PersistenceService struct {
	producer kafka.Producer
	consumer kafka.Consumer
	db       db.Connection
}

func New(producer kafka.Producer, consumer kafka.Consumer, db db.Connection) Service {
	return &PersistenceService{
		producer: producer,
		consumer: consumer,
		db:       db,
	}
}

func (s *PersistenceService) Consume(ctx context.Context, topic string) error {
	return s.consumer.Consume(ctx, []string{topic}, model.NewConsumer(s.db))
}

func (s *PersistenceService) Produce(ctx context.Context, msg *sarama.ProducerMessage) error {
	var wg sync.WaitGroup
	wg.Add(1)

	var err error

	go func() {
		defer wg.Done()

		select {
		case <-s.producer.Successes():
			slog.Info("message delivery acknowledged", "message", msg)
		case err = <-s.producer.Errors():
			slog.Error("failed to produce message", "error", err, "message", msg)
		case <-ctx.Done():
			slog.Info("producer closed", "message", msg)
		}
	}()

	s.producer.Input() <- msg

	wg.Wait()

	return err
}
