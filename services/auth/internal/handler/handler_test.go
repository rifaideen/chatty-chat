package handler

import (
	mock_service "auth/internal/service/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pkg/auth"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockService(ctrl)

	handler := New(service)

	t.Run("login error", func(t *testing.T) {
		// create test http request with body for login handler
		request := &auth.LoginRequest{
			Username: "test",
			Password: "test",
		}

		w := httptest.NewRecorder()

		payload, _ := json.Marshal(request)

		// create test http request
		r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payload))

		service.EXPECT().Login(gomock.Any()).DoAndReturn(func(request *auth.LoginRequest) (string, error) {
			return "", errors.New("invalid credentials")
		})

		// call login handler
		handler.Login(w, r)

		// assert failure
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("login success", func(t *testing.T) {
		// create test http request with body for login handler
		request := &auth.LoginRequest{
			Username: "admin",
			Password: "admin",
		}

		w := httptest.NewRecorder()

		payload, _ := json.Marshal(request)

		// create test http request
		r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payload))
		service.EXPECT().Login(gomock.Any()).DoAndReturn(func(request *auth.LoginRequest) (string, error) {
			return "token", nil
		})

		// call login handler
		handler.Login(w, r)

		// assert success
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestVerifyHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockService(ctrl)

	handler := New(service)

	t.Run("verify error", func(t *testing.T) {
		// create test http request with body for login handler
		request := &auth.VerifyRequest{
			Token: "token",
		}

		w := httptest.NewRecorder()

		payload, _ := json.Marshal(request)

		// create test http request
		r := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer(payload))

		service.EXPECT().VerifyToken(gomock.Any()).DoAndReturn(func(request *auth.VerifyRequest) (*jwt.Token, error) {
			return nil, errors.New("invalid token")
		})
		// call login handler
		handler.Verify(w, r)
		// assert failure
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("verify success", func(t *testing.T) {
		// create test http request with body for login handler
		request := &auth.VerifyRequest{
			Token: "token",
		}

		w := httptest.NewRecorder()

		payload, _ := json.Marshal(request)

		// create test http request
		r := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer(payload))

		service.EXPECT().VerifyToken(gomock.Any()).DoAndReturn(func(request *auth.VerifyRequest) (*jwt.Token, error) {
			return &jwt.Token{}, nil
		})

		// call verify handler
		handler.Verify(w, r)

		// assert success
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
