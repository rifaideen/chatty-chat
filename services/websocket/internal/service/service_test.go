package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pkg/auth"
	mock_kafka "pkg/kafka/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_producer := mock_kafka.NewMockProducer(ctrl)
	mock_consumer := mock_kafka.NewMockConsumer(ctrl)
	topics := []string{"test-consumer", "test-producer"}

	authServiceUrl := "http://localhost:8080"
	ollamaServiceUrl := "http://localhost:8081"

	service := New(mock_consumer, mock_producer, topics, authServiceUrl, ollamaServiceUrl)

	t.Run("context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go service.Listen(ctx)

		time.Sleep(time.Second)
	})
}

func TestVerify(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_producer := mock_kafka.NewMockProducer(ctrl)
	mock_consumer := mock_kafka.NewMockConsumer(ctrl)
	topics := []string{"test-consumer", "test-producer"}

	// Create a test server that simulates the external service
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request body is as expected
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path == "/auth/verify" {
			req := auth.VerifyRequest{}
			err := json.NewDecoder(r.Body).Decode(&req)

			if err != nil {
				t.Errorf("Error decoding request body: %v", err)
				return
			}

			w.WriteHeader(http.StatusOK)

			var response *auth.VerifyResponse

			if req.Token == "test-token" {
				response = &auth.VerifyResponse{
					Valid: true,
				}
			} else {
				response = &auth.VerifyResponse{
					Valid: false,
					Error: "Invalid token",
				}
			}

			json.NewEncoder(w).Encode(response)

			return
		}

		// ... check other request details (headers, URL, body) ...

		// Simulate a successful response
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	authServiceUrl := ts.URL
	ollamaServiceUrl := ts.URL

	service := New(mock_consumer, mock_producer, topics, authServiceUrl, ollamaServiceUrl)

	t.Run("verify success", func(t *testing.T) {
		verified, err := service.Verify("test-token")
		assert.True(t, verified)
		assert.NoError(t, err)
	})

	t.Run("verify error", func(t *testing.T) {
		verified, err := service.Verify("invalid-token")

		assert.False(t, verified)
		assert.Error(t, err)
	})
}

func TestConsume(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_producer := mock_kafka.NewMockProducer(ctrl)
	mock_consumer := mock_kafka.NewMockConsumer(ctrl)
	topics := []string{"test-consumer", "test-producer"}
	authServiceUrl := "http://localhost:8080"
	ollamaServiceUrl := "http://localhost:8081"

	service := New(mock_consumer, mock_producer, topics, authServiceUrl, ollamaServiceUrl)

	t.Run("consume success", func(t *testing.T) {
		mock_consumer.EXPECT().Consume(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		err := service.Consume()

		assert.NoError(t, err)
	})

	t.Run("consume error", func(t *testing.T) {
		mock_consumer.EXPECT().Consume(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
		err := service.Consume()
		assert.Error(t, err)
	})
}
