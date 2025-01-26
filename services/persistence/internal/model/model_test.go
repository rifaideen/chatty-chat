package model

import (
	mock_db "pkg/db/mocks"
	"testing"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_db := mock_db.NewMockConnection(ctrl)

	consumer := NewConsumer(mock_db)
	assert.NotNil(t, consumer)

	err := consumer.Setup(nil)
	assert.NoError(t, err)
}

func TestConsumeClaim(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_db := mock_db.NewMockConnection(ctrl)

	consumer := NewConsumer(mock_db)
	assert.NotNil(t, consumer)

	err := consumer.ConsumeClaim(nil, nil)
	assert.Error(t, err)
}

func TestCleanup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_db := mock_db.NewMockConnection(ctrl)

	consumer := NewConsumer(mock_db)
	assert.NotNil(t, consumer)

	err := consumer.Cleanup(nil)
	assert.NoError(t, err)
}

func TestStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_db := mock_db.NewMockConnection(ctrl)
	consumer := NewConsumer(mock_db)

	assert.NotNil(t, consumer)

	t.Run("store error", func(t *testing.T) {
		err := consumer.store(nil)

		assert.EqualError(t, err, "message is nil")
	})

	t.Run("store success", func(t *testing.T) {
		query := `
        INSERT INTO messages (
            message,
            topic,
            partition,
            ofset,
            timestamp,
            created_at
        ) VALUES ($1, $2, $3, $4, $5, NOW())
    `

		mock_db.EXPECT().Exec(query, gomock.Any()).Return(nil, nil)
		err := consumer.store(&sarama.ConsumerMessage{
			Topic:     "test",
			Partition: 1,
			Offset:    1,
			Value:     []byte("test"),
		})

		assert.NoError(t, err)
	})
}
