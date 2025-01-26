![](assets/20250125_191502_logo.svg)

AI powered chat application built with Vue, Golang, WebSocket, Postgres, Kafka and Docker.

## Features

- real-time chatting
- authentication
- real-time notifications
- docker support for easy deployment

## Tech Stack

- Vue.js
- Golang
- Postgres
- Kafka
- Docker

## Prerequisites

- Docker and Docker Compose installed on your system
- Git

## Quick Start

1. Clone the repository:

```bash
git clone https://github.com/rifaideen/chatty-chat.git
```

2.Start the application using Docker:

```bash
docker compose up --build -d
```

The application will be available at `http://localhost:8004`

# Microservices

- **Auth Microservice:** http://localhost:8001
- **Persistence Microservice:** (background service)
- **Chat Microservice (Websocket):** http://localhost:8003
- **UI Microservice:** http://localhost:8004

# Application usage:

1. **Login:** credentials `admin` and `admin`
2. **Pull:** navigate to `Models` page and pull any model
3. **Chat:** chat with AI

# Create SSL Certificate
Please create ssl certificates
`openssl req -x509 -newkey rsa:2048 -nodes -keyout certs/server.key -out certs/server.crt -days 365`

Change the permission
`chmod 600 certs/server.key`
