package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"pkg/ai"
	"pkg/auth"
	"pkg/kafka"
	"sync"

	"github.com/IBM/sarama"
	"github.com/gorilla/websocket"
	"github.com/rifaideen/talkative"
)

type WebsocketService interface {
	Consume() error
	Listen(ctx context.Context)
	ServeWS(w http.ResponseWriter, r *http.Request)
	Verify(token string) (bool, error)
}

// Service maintains the set of active clients and broadcasts messages to them.
type Service struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex

	producer       sarama.AsyncProducer
	consumer       sarama.ConsumerGroup
	topicConsumer  string
	topicProducer  string
	authServiceUrl string
	ollama         *talkative.Client
	ai             *ai.AI
}

// New initializes and returns a new Service.
func New(consumer kafka.Consumer, producer kafka.Producer, topics []string, authServiceUrl, ollamaServiceUrl string) WebsocketService {
	client, err := talkative.New(ollamaServiceUrl)

	if err != nil {
		log.Fatal(err)
	}

	ai, err := ai.New(ollamaServiceUrl)

	if err != nil {
		log.Fatal(err)
	}

	return &Service{
		clients:        make(map[*Client]bool),
		broadcast:      make(chan []byte),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		producer:       producer,
		consumer:       consumer,
		topicProducer:  topics[0],
		topicConsumer:  topics[1],
		authServiceUrl: authServiceUrl,
		ollama:         client,
		ai:             ai,
	}
}

// Listen starts the service to handle client registration, unregistration, and broadcasting messages.
func (m *Service) Listen(ctx context.Context) {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client] = true
			m.mu.Unlock()

			log.Println("New client connected")

		case client := <-m.unregister:
			m.mu.Lock()

			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client.send)
			}

			m.mu.Unlock()
			log.Println("Client disconnected")

		case message := <-m.broadcast:
			m.mu.Lock()

			for client := range m.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(m.clients, client)
				}
			}

			m.mu.Unlock()
		case <-ctx.Done():
			log.Println("Context cancelled, stopping service")
			return
		}
	}
}

func (m *Service) ServeWS(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	// verify the token
	token := r.URL.Query().Get("token")
	valid, err := m.Verify(token)

	if err != nil {
		conn.WriteJSON(map[string]string{
			"error": "unable to verify token",
		})
		conn.Close()

		log.Println("Error verifying token:", err)
		return
	}

	if !valid {
		conn.WriteJSON(map[string]string{
			"error": "Invalid token",
		})
		conn.Close()

		log.Println("Invalid token")
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	m.register <- client

	// start the read and write goroutines
	go client.read(m)
	go client.write()
}

// verify the token with auth service and return the result
func (m *Service) Verify(token string) (bool, error) {
	data, err := json.Marshal(auth.VerifyRequest{
		Token: token,
	})

	if err != nil {
		return false, err
	}

	url := fmt.Sprintf("%s/auth/verify", m.authServiceUrl)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))

	if err != nil {
		return false, err
	}

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return false, nil
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return false, err
	}

	verification := auth.VerifyResponse{}

	err = json.Unmarshal(body, &verification)

	if err != nil {
		return false, err
	}

	if verification.Error != "" {
		return false, fmt.Errorf("%s", verification.Error)
	}

	return verification.Valid, nil
}

func (m *Service) Consume() error {
	slog.Info("cosuming messages from kafka", "topic", m.topicConsumer)

	return m.consumer.Consume(context.Background(), []string{m.topicConsumer}, &Consumer{
		manager: m,
	})
}
