package service

import (
	"encoding/json"
	"log/slog"

	"github.com/IBM/sarama"
)

// Consumer handles Kafka message consumption and broadcasts to websocket clients
type Consumer struct {
	manager *Service // Reference to websocket connection manager service
}

// Setup is called when the consumer group session starts
func (*Consumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is called when the consumer group session ends
func (*Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from a Kafka partition
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			// Check if message channel is still open
			if !ok {
				slog.Info("message channel closed. exiting...")
				return nil
			}

			message := string(msg.Value)

			// Log received message details
			slog.Info("message received from kafka", "topic", msg.Topic, "message", message)

			data, err := json.Marshal(map[string]interface{}{
				"type": "notification",
				"data": message,
			})

			if err != nil {
				slog.Error("error marshalling message", "error", err)
				continue
			}

			// Broadcast message to all connected websocket clients
			c.manager.broadcast <- data

			// Mark message as processed
			session.MarkMessage(msg, "")
		case <-session.Context().Done():
			// Exit when consumer group session is done
			return nil
		}
	}
}
