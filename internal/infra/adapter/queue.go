package adapter

import (
	"context"
	"encoding/json"
	"os"

	"github.com/marcoscoutinhodev/mv_chat/internal/entity"
	"github.com/marcoscoutinhodev/mv_chat/internal/infra/amqp"
)

type Queue struct {
	amqp                         amqp.AMQP
	registerNotificationQueue    string
	registerNotificationExchange string
}

func NewQueue() *Queue {
	return &Queue{
		amqp:                         *amqp.NewAMQP(),
		registerNotificationQueue:    os.Getenv("REGISTER_NOTIFICATION_QUEUE"),
		registerNotificationExchange: os.Getenv("REGISTER_NOTIFICATION_EXCHANGE"),
	}
}

func (q Queue) RegisterNotification(ctx context.Context, user *entity.User) error {
	input := map[string]string{
		"name":  user.Name,
		"email": user.Email,
	}

	body, err := json.Marshal(&input)
	if err != nil {
		return err
	}

	if err = q.amqp.Producer(ctx, q.registerNotificationQueue, q.registerNotificationExchange, body); err != nil {
		return err
	}

	return nil
}
