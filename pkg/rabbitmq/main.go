package rabbitmq

import (
	"context"
	"encoding/json"
	"scheduler/internal/config"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQ(config *config.RabbitMQConfig) (*RabbitMQ, error) {
	conn, err := amqp.Dial(config.Url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rmq *RabbitMQ) ensureQueue(name string, durable bool, autoDelete bool) (*amqp.Queue, error) {
	q, err := rmq.ch.QueueDeclare(strings.ToLower(name), durable, autoDelete, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (rmq *RabbitMQ) AddDurableQueue(name, exchange, routingKey string) error {
	q, err := rmq.ensureQueue(name, true, false)
	if err != nil {
		return err
	}

	err = rmq.ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = rmq.ch.QueueBind(q.Name, routingKey, exchange, false, nil)

	return err
}

func (rmq *RabbitMQ) Publish(key string, exchangeName string, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rmq.ch.PublishWithContext(ctx, exchangeName, key, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         data,
		Timestamp:    time.Now(),
	})

	return err
}

func (rmq *RabbitMQ) Consume(
	ctx context.Context,
	queue string,
	handler func(any) error,
) error {
	msgs, err := rmq.ch.Consume(queue, "scheduler", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case msg := <-msgs:
			var data any

			if err := json.Unmarshal(msg.Body, &data); err != nil {
				_ = msg.Nack(false, true)
				continue
			}

			if err := handler(data); err != nil {
				_ = msg.Nack(false, true)
				continue
			}

			_ = msg.Ack(false)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
