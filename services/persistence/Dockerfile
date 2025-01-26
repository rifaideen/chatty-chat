FROM golang:1.23-alpine
 
WORKDIR /app
 
# Copy go mod and sum files
COPY ./services/persistence/go.mod ./services/persistence/go.sum ./
 
# Download all dependencies
RUN go mod download
 
# Copy the source code
COPY . .
 
# Build the application
RUN go build -o main ./services/persistence/cmd
 
# Expose port 8002
EXPOSE 8002
 
# Command to run the executable
CMD ["./main"]