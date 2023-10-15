package infra

import (
	"time"

	"github.com/marcoscoutinhodev/ms_email_notification/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *amqp.Connection
}

func NewConsumer() *Consumer {
	consumer := &Consumer{}
	consumer.connect()
	return consumer
}

func (c *Consumer) connect() {
	if c.conn == nil || c.conn.IsClosed() {
		newConn, err := amqp.Dial(config.RBMQ_URI)
		if err != nil {
			time.Sleep(time.Second * 30)
			c.connect()
			return
		}

		c.conn = newConn
	}
}

func (c Consumer) Delivery() (<-chan amqp.Delivery, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		c.connect()
		return c.Delivery()
	}

	return ch.Consume(
		config.EMAIL_NOTIFICATION_QUEUE,
		*new(string),
		false,
		false,
		false,
		false,
		nil,
	)
}

func (c Consumer) Close() {
	c.conn.Close()
}
