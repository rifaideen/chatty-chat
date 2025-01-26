package model

import (
	"errors"
	"log"
	"pkg/db"

	"github.com/IBM/sarama"
)

type Consumer struct {
	db db.Connection
}

func NewConsumer(db db.Connection) *Consumer {
	return &Consumer{
		db: db,
	}
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer setup completed.")

	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	if session == nil || claim == nil {
		return errors.New("session or claim cannot be nil")
	}

	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				log.Println("Message channel closed. Exiting...")
				return nil
			}

			const maxRetries = 5
			var err error

			for attempt := 0; attempt < maxRetries; attempt++ {
				err = c.store(msg)

				if err == nil {
					session.MarkMessage(msg, "")
					break
				}

				log.Printf("Retry attempt %d/%d failed: %v", attempt+1, maxRetries, err)

				if attempt == maxRetries-1 {
					log.Printf("Retry attempt permenantly failed after %d attempts", maxRetries)
					session.MarkMessage(msg, "")
				}
			}
		case <-session.Context().Done():
			return nil
		}

	}
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer cleanup completed.")

	return nil
}

func (c *Consumer) store(msg *sarama.ConsumerMessage) error {
	if msg == nil {
		return errors.New("message is nil")
	}

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
	_, err := c.db.Exec(
		query,
		string(msg.Value),
		msg.Topic,
		msg.Partition,
		msg.Offset,
		msg.Timestamp,
	)

	return err
}
