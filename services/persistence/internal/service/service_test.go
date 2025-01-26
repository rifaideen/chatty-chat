package service

import (
	"context"
	mock_db "pkg/db/mocks"
	mock_kafka "pkg/kafka/mocks"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestConsume(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_producer := mock_kafka.NewMockProducer(ctrl)
	mock_consumer := mock_kafka.NewMockConsumer(ctrl)
	mock_db := mock_db.NewMockConnection(ctrl)

	service := New(mock_producer, mock_consumer, mock_db)

	t.Run("consume success", func(t *testing.T) {
		mock_consumer.EXPECT().Consume(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		err := service.Consume(context.Background(), "test")

		assert.NoError(t, err)
	})

	t.Run("consume failure", func(t *testing.T) {
		mock_consumer.EXPECT().Consume(gomock.Any(), gomock.Any(), gomock.Any()).Return(assert.AnError)
		err := service.Consume(context.Background(), "test")
		assert.Error(t, err)
	})
}

func TestProduce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_producer := mock_kafka.NewMockProducer(ctrl)
	mock_consumer := mock_kafka.NewMockConsumer(ctrl)
	mock_db := mock_db.NewMockConnection(ctrl)

	service := New(mock_producer, mock_consumer, mock_db)

	t.Run("produce context timeout", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			defer cancel()
			time.Sleep(time.Millisecond * 500)
		}()

		mock_producer.EXPECT().Successes().Return(make(chan *sarama.ProducerMessage))
		mock_producer.EXPECT().Errors().Return(make(chan *sarama.ProducerError))
		mock_producer.EXPECT().Input().DoAndReturn(func() chan<- *sarama.ProducerMessage {
			return make(chan *sarama.ProducerMessage, 1)
		})
		err := service.Produce(ctx, &sarama.ProducerMessage{})
		assert.NoError(t, err)
	})

	t.Run("produce success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			defer cancel()
			time.Sleep(time.Second * 2)
		}()

		mock_producer.EXPECT().Successes().DoAndReturn(func() <-chan *sarama.ProducerMessage {
			msg := make(chan *sarama.ProducerMessage, 1)
			msg <- &sarama.ProducerMessage{}

			return msg
		})
		mock_producer.EXPECT().Errors().Return(make(chan *sarama.ProducerError))
		mock_producer.EXPECT().Input().Return(make(chan *sarama.ProducerMessage, 1))

		err := service.Produce(ctx, &sarama.ProducerMessage{})
		assert.NoError(t, err)
	})

	t.Run("produce failure", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			defer cancel()
			time.Sleep(time.Second * 2)
		}()

		mock_producer.EXPECT().Successes().Return(make(chan *sarama.ProducerMessage))
		mock_producer.EXPECT().Errors().DoAndReturn(func() <-chan *sarama.ProducerError {
			err := make(chan *sarama.ProducerError, 1)
			err <- &sarama.ProducerError{Err: assert.AnError}

			return err
		})
		mock_producer.EXPECT().Input().Return(make(chan *sarama.ProducerMessage, 1))

		err := service.Produce(ctx, &sarama.ProducerMessage{})
		assert.Error(t, err)
	})
}
