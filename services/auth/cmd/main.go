package main

import (
	"auth/internal/config"
	"auth/internal/handler"
	"auth/internal/service"
	"log"
	"log/slog"
	"net/http"
	"os"
	"pkg/middleware"
)

const PORT = ":8001"

func main() {
	handler := bootstrap()

	http.HandleFunc("POST /auth/login", handler.Login)
	http.HandleFunc("POST /auth/verify", handler.Verify)

	log.Println("auth service listening on http://localhost" + PORT)

	log.Fatal(http.ListenAndServe(PORT, middleware.Cors(http.DefaultServeMux)))
}

func bootstrap() handler.Handler {
	config := config.Load()

	// configure logger with JSON formatter and set as default
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	slog.SetDefault(log)

	service := service.New(
		config.Secret,
		config.ExpiresIn,
	)

	return handler.New(service)
}
