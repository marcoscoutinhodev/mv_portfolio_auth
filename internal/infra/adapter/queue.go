package adapter

import (
	"context"
	"encoding/json"
	"os"

	"github.com/marcoscoutinhodev/mv_chat/internal/entity"
	"github.com/marcoscoutinhodev/mv_chat/internal/infra/amqp"
)

type Queue struct {
	amqp                                  amqp.AMQP
	registerNotificationKey               string
	registerNotificationQueue             string
	registerNotificationExchange          string
	forgottenPasswordNotificationKey      string
	forgottenPasswordNotificationQueue    string
	forgottenPasswordNotificationExchange string
}

func NewQueue() *Queue {
	return &Queue{
		amqp:                                  *amqp.NewAMQP(),
		registerNotificationQueue:             os.Getenv("REGISTER_NOTIFICATION_QUEUE"),
		registerNotificationKey:               os.Getenv("REGISTER_NOTIFICATION_KEY"),
		registerNotificationExchange:          os.Getenv("REGISTER_NOTIFICATION_EXCHANGE"),
		forgottenPasswordNotificationKey:      os.Getenv("FORGOT_PASSWORD_NOTIFICATION_KEY"),
		forgottenPasswordNotificationQueue:    os.Getenv("FORGOT_PASSWORD_NOTIFICATION_QUEUE"),
		forgottenPasswordNotificationExchange: os.Getenv("FORGOT_PASSWORD_NOTIFICATION_EXCHANGE"),
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

	if err = q.amqp.Producer(ctx, q.registerNotificationKey, q.registerNotificationExchange, body); err != nil {
		return err
	}

	return nil
}

func (q Queue) ForgottenPasswordNotification(ctx context.Context, user *entity.User, token string) error {
	input := map[string]string{
		"name":  user.Name,
		"email": user.Email,
		"token": token,
	}

	body, err := json.Marshal(&input)
	if err != nil {
		return err
	}

	if err = q.amqp.Producer(ctx, q.forgottenPasswordNotificationKey, q.forgottenPasswordNotificationExchange, body); err != nil {
		return err
	}

	return nil
}
