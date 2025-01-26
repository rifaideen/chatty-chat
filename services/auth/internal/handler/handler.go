package handler

import (
	"auth/internal/service"
	"encoding/json"
	"log/slog"
	"net/http"
	"pkg/auth"
	"strings"
)

type Handler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Verify(w http.ResponseWriter, r *http.Request)
}

type AuthHandler struct {
	service service.Service
}

func New(service service.Service) Handler {
	return &AuthHandler{
		service: service,
	}
}

func (c *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	slog.Info("login request received")

	// set response header
	w.Header().Add("Content-Type", "application/json")

	writer := json.NewEncoder(w)

	request := &auth.LoginRequest{}

	json.NewDecoder(r.Body).Decode(request)

	token, err := c.service.Login(request)

	if err != nil {
		slog.Warn("login failed", "error", err.Error(), "request", request)

		w.WriteHeader(http.StatusBadRequest)
		writer.Encode(auth.LoginResponse{
			Error: err.Error(),
		})

		return
	}

	w.WriteHeader(200)

	writer.Encode(auth.LoginResponse{
		Token: token,
		User:  strings.Title(request.Username),
	})

	slog.Info("login successful")
}

func (c *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	slog.Info("token verification request received")

	// set response header
	w.Header().Add("Content-Type", "application/json")

	writer := json.NewEncoder(w)

	request := &auth.VerifyRequest{}
	json.NewDecoder(r.Body).Decode(request)

	token, err := c.service.VerifyToken(request)

	if err != nil {
		slog.Warn("token verification failed", "error", err.Error())

		w.WriteHeader(http.StatusBadRequest)
		writer.Encode(auth.VerifyResponse{
			Error: "invalid token",
		})
		return
	}

	w.WriteHeader(200)

	writer.Encode(auth.VerifyResponse{
		Valid: true,
	})

	slog.Info("token verification successful", "claims", token.Claims)
}
