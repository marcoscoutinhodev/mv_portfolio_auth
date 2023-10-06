package adapter

import (
	"context"
	"encoding/json"

	"github.com/marcoscoutinhodev/ms_auth/config"
	"github.com/marcoscoutinhodev/ms_auth/internal/entity"
	"github.com/marcoscoutinhodev/ms_auth/internal/infra/amqp"
)

type EmailNotification struct {
	amqp                      amqp.AMQP
	registerKey               string
	registerExchange          string
	forgottenPasswordKey      string
	forgottenPasswordExchange string
}

func NewEmailNotification() *EmailNotification {
	return &EmailNotification{
		amqp:                      *amqp.NewAMQP(),
		registerKey:               config.REGISTER_NOTIFICATION_KEY,
		registerExchange:          config.REGISTER_NOTIFICATION_EXCHANGE,
		forgottenPasswordKey:      config.FORGOT_PASSWORD_NOTIFICATION_KEY,
		forgottenPasswordExchange: config.FORGOT_PASSWORD_NOTIFICATION_EXCHANGE,
	}
}

func (q EmailNotification) producer(ctx context.Context, key, exchange string, input interface{}) error {
	body, err := json.Marshal(&input)
	if err != nil {
		return err
	}

	if err = q.amqp.Producer(ctx, key, exchange, body); err != nil {
		return err
	}

	return nil
}

func (q EmailNotification) Register(ctx context.Context, user *entity.User, token string) error {
	input := map[string]string{
		"name":  user.Name,
		"email": user.Email,
		"token": token,
	}

	if err := q.producer(ctx, q.registerKey, q.registerExchange, &input); err != nil {
		return err
	}

	return nil
}

func (q EmailNotification) ForgottenPassword(ctx context.Context, user *entity.User, token string) error {
	input := map[string]string{
		"name":  user.Name,
		"email": user.Email,
		"token": token,
	}

	if err := q.producer(ctx, q.forgottenPasswordKey, q.forgottenPasswordExchange, &input); err != nil {
		return err
	}

	return nil
}
