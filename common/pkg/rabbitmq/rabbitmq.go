package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQ creates a reusable RabbitMQ instance
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil
}

// DeclareExchange sets up an exchange (can be called multiple times)
func (r *RabbitMQ) DeclareExchange(name, kind string, durable bool) error {
	return r.channel.ExchangeDeclare(
		name, kind, durable, false, false, false, nil,
	)
}

// DeclareQueue sets up a queue (can be used for multiple queues)
func (r *RabbitMQ) DeclareQueue(name string, durable bool) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		name, durable, false, false, false, nil,
	)
}

// BindQueue binds a queue to an exchange with a routing key
func (r *RabbitMQ) BindQueue(queueName, exchangeName, routingKey string) error {
	return r.channel.QueueBind(
		queueName, routingKey, exchangeName, false, nil,
	)
}

// Publish sends a message to a specific exchange with a routing key
func (r *RabbitMQ) Publish(exchange, routingKey string, body []byte, contentType string) error {
	return r.channel.Publish(
		exchange,
		routingKey,
		true,
		false,
		amqp.Publishing{
			ContentType:  contentType,
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

// Consume returns a delivery channel for the specified queue
func (r *RabbitMQ) Consume(queueName string, autoAck bool) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queueName, "", autoAck, false, false, false, nil,
	)
}

// Close gracefully closes the channel and connection
func (r *RabbitMQ) Close() {
	if r.channel != nil {
		_ = r.channel.Close()
	}
	if r.conn != nil {
		_ = r.conn.Close()
	}
}

func (r *RabbitMQ) Build(queueName, exchange, routingKey string) error {
	err := r.DeclareExchange(exchange, "direct", true)
	if err != nil {
		return fmt.Errorf("exchange declare failed: %w", err)
	}

	queue, err := r.DeclareQueue(queueName, true)
	if err != nil {
		return fmt.Errorf("queue declare failed: %w", err)
	}

	err = r.BindQueue(queue.Name, exchange, routingKey)
	if err != nil {
		return fmt.Errorf("queue bind failed: %w", err)
	}

	return nil
}
