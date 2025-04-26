FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY src/ ./src/    
RUN go build -o main ./src/cmd/main.go

EXPOSE 9090

CMD ["./main"]