package queues

import (
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type NotificationsQueue struct {
	channel *amqp.Channel
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
	}, nil
}
