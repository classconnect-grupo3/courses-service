package queues

import (
	"encoding/json"
	"fmt"

	"courses-service/src/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type NotificationsQueueInterface interface {
	Publish(message QueueMessage) error
}

type NotificationsQueue struct {
	channel   *amqp.Channel
	queueName string
}

func NewNotificationsQueue(config *config.Config) (*NotificationsQueue, error) {
	queueName := config.NotificationsQueueName
	queueURL := config.RabbitMQURL

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
		channel:   ch,
		queueName: queueName,
	}, nil
}

func (q *NotificationsQueue) Publish(message QueueMessage) error {
	if q.channel == nil {
		return nil // testing purposes
	}

	body, err := message.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return q.channel.Publish(
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
