FROM golang:1.23-alpine
 
WORKDIR /app
 
# Copy go mod and sum files
COPY ./services/auth/go.mod ./services/auth/go.sum ./
 
# Download all dependencies
RUN go mod download
 
# Copy the source code
COPY . .
 
# Build the application
RUN go build -o main ./services/auth/cmd
 
# Expose port 8001
EXPOSE 8001
 
# Command to run the executable
CMD ["./main"]