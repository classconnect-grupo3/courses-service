FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY src/ ./src/
RUN go build -o main ./src/main.go

EXPOSE 9090

CMD ["./main"]