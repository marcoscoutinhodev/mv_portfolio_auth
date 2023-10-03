package adapter

import (
	"context"
	"encoding/json"
	"os"

	"github.com/marcoscoutinhodev/mv_chat/internal/entity"
	"github.com/marcoscoutinhodev/mv_chat/internal/infra/amqp"
)

type EmailNotification struct {
	amqp                               amqp.AMQP
	registerKey                        string
	registerEmailNotification          string
	registerExchange                   string
	forgottenPasswordKey               string
	forgottenPasswordEmailNotification string
	forgottenPasswordExchange          string
}

func NewEmailNotification() *EmailNotification {
	return &EmailNotification{
		amqp:                               *amqp.NewAMQP(),
		registerEmailNotification:          os.Getenv("REGISTER_NOTIFICATION_QUEUE"),
		registerKey:                        os.Getenv("REGISTER_NOTIFICATION_KEY"),
		registerExchange:                   os.Getenv("REGISTER_NOTIFICATION_EXCHANGE"),
		forgottenPasswordKey:               os.Getenv("FORGOT_PASSWORD_NOTIFICATION_KEY"),
		forgottenPasswordEmailNotification: os.Getenv("FORGOT_PASSWORD_NOTIFICATION_QUEUE"),
		forgottenPasswordExchange:          os.Getenv("FORGOT_PASSWORD_NOTIFICATION_EXCHANGE"),
	}
}

func (q EmailNotification) Register(ctx context.Context, user *entity.User) error {
	input := map[string]string{
		"name":  user.Name,
		"email": user.Email,
	}

	body, err := json.Marshal(&input)
	if err != nil {
		return err
	}

	if err = q.amqp.Producer(ctx, q.registerKey, q.registerExchange, body); err != nil {
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

	body, err := json.Marshal(&input)
	if err != nil {
		return err
	}

	if err = q.amqp.Producer(ctx, q.forgottenPasswordKey, q.forgottenPasswordExchange, body); err != nil {
		return err
	}

	return nil
}
