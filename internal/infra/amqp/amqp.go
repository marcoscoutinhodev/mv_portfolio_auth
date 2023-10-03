package amqp

import (
	"context"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQP struct {
	conn *amqp.Connection
}

func NewAMQP() *AMQP {
	c, err := amqp.Dial(os.Getenv("RBMQ_URI"))
	if err != nil {
		panic(err)
	}

	return &AMQP{
		conn: c,
	}
}

func (a AMQP) Producer(ctx context.Context, emailNotification, exchange string, body []byte) error {
	if a.conn.IsClosed() {
		c, err := amqp.Dial(os.Getenv("RBMQ_URI"))
		if err != nil {
			return err
		}

		a.conn = c
	}

	ch, err := a.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	if err := ch.PublishWithContext(
		ctx,
		exchange,
		emailNotification,
		true,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	); err != nil {
		return err
	}

	return nil
}
