package queues

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type NotificationsQueue struct {
	channel *amqp.Channel
	queueName string
}

func NewNotificationsQueue() (*NotificationsQueue, error) {
	queueName := os.Getenv("NOTIFICATIONS_QUEUE_NAME")
	queueURL := os.Getenv("RABBITMQ_URL")

	conn, err := amqp.Dial(queueURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &NotificationsQueue{
		channel: ch,
		queueName: queueName,
	}, nil
}

func (q *NotificationsQueue) Publish(message QueueMessage) error {
	body, err := message.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return q.channel.PublishWithContext(
		context.Background(),
		"",
		q.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBody,
		},
	)
}
