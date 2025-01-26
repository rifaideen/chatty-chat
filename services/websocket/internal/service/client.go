package service

import (
	"encoding/json"
	"log"
	"log/slog"
	"pkg/ai"
	"websocket/internal/model"

	"github.com/IBM/sarama"
	"github.com/gorilla/websocket"
	"github.com/rifaideen/talkative"
)

// Client represents a single WebSocket connection with a send channel for messages
type Client struct {
	// conn holds the WebSocket connection instance
	conn *websocket.Conn
	// send is a channel for buffering outbound messages
	send chan []byte
}

// read handles incoming messages from the WebSocket connection
func (c *Client) read(manager *Service) {
	defer func() {
		manager.unregister <- c
		c.conn.Close()
	}()

	go func() {
		for {
			select {
			case _, ok := <-manager.producer.Successes():
				if !ok {
					slog.Info("producer channel closed")

					return
				}

				slog.Info("message sent successfully")
			case err := <-manager.producer.Errors():
				slog.Error("error sending message", "error", err)

				return
			}
		}
	}()

	for {
		message := &model.Message{}
		err := c.conn.ReadJSON(message)

		if err != nil {
			slog.Error("error reading message", "error", err)
			break
		}

		slog.Info("received message from client", "chat", message.Data)

		// forward the message to the producer topic in kafka and then initiate a chat
		manager.producer.Input() <- &sarama.ProducerMessage{
			Topic: manager.topicProducer,
			Value: sarama.StringEncoder(message.Data),
		}

		if message.Type == "pull" {
			// pull the AI model
			c.pull(manager, message.Model)
		} else {
			// chat with the AI
			c.chat(manager, message.Model, message.Data)
		}
	}
}

// write continuously listens on the send channel and writes messages to the WebSocket connection
// It handles the outbound message flow until an error occurs or the connection closes
func (c *Client) write() {
	defer c.conn.Close()

	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			slog.Error("error writing message", "error", err, "message", string(message))
			break
		}
	}
}

func (c *Client) chat(manager *Service, model, prompt string) {
	// Callback function to handle the response
	callback := func(cr *talkative.ChatResponse, err error) {
		if err != nil {
			slog.Error("unable to process chat response", "error", err)
			return
		}

		data, err := json.Marshal(map[string]interface{}{
			"type": "chat",
			"data": cr.Message.Content,
			"done": cr.Done,
		})

		if err != nil {
			slog.Error("unable to marshal json respose", "error", err)
			return
		}

		c.send <- data
	}
	// Additional parameters to include. (Optional)
	var params *talkative.ChatParams = nil
	// The chat message to send
	message := talkative.ChatMessage{
		Role:    talkative.USER, // Initiate the chat as a user
		Content: prompt,
	}

	done, err := manager.ollama.Chat(model, callback, params, message)

	if err != nil {
		panic(err)
	}

	<-done // wait for the chat to complete
}

func (c *Client) pull(manager *Service, model string) {
	// Callback function to handle the response
	callback := func(cr *ai.PullResponse, err error) {
		if err != nil {
			slog.Error("unable to process pull response", "error", err)
			return
		}

		data, err := json.Marshal(map[string]interface{}{
			"type": "pull",
			"data": cr,
			"done": cr.Status == "success" || cr.Status == "writing manifest",
		})

		if err != nil {
			slog.Error("unable to marshal json respose", "error", err)
			return
		}

		c.send <- data
	}

	done, err := manager.ai.Pull(model, callback)

	if err != nil {
		panic(err)
	}

	<-done // wait for the chat to complete
}

// Setup implements the ConsumerGroupHandler interface
// Called when the consumer group session is set up
func (*Client) Setup(_ sarama.ConsumerGroupSession) error {
	log.Println("Consumer setup completed.")
	return nil
}

// Cleanup implements the ConsumerGroupHandler interface
// Called when the consumer group session is cleaning up
func (*Client) Cleanup(_ sarama.ConsumerGroupSession) error {
	log.Println("Consumer cleanup completed.")
	return nil
}

// ConsumeClaim implements the ConsumerGroupHandler interface
// Handles the consumption of messages from a Kafka topic partition
// Forwards received messages to the WebSocket connection via the send channel
func (c *Client) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				log.Println("Message channel closed. Exiting...")
				return nil
			}

			log.Printf("Message received: topic=%s partition=%d offset=%d value=%s",
				msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
			c.send <- msg.Value
			session.MarkMessage(msg, "")
		case <-session.Context().Done():
			return nil
		}

	}
}
