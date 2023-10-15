package service

import (
	"context"
	"encoding/json"

	"github.com/marcoscoutinhodev/ms_email_notification/internal/entity"
)

type Email struct {
	emailProvider EmailProviderInterface
}

func NewEmail(emailProvider EmailProviderInterface) *Email {
	return &Email{
		emailProvider: emailProvider,
	}
}

type AuthNotificationInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Title    string `json:"title"`
	Template string `json:"template"`
}

func (e Email) AuthNotification(ctx context.Context, in []byte) error {
	var input AuthNotificationInput
	if err := json.Unmarshal(in, &input); err != nil {
		return err
	}

	buff, err := TemplateGenerator(input)
	if err != nil {
		return err
	}

	recipient := entity.NewRecipient(input.Name, input.Email)
	mail := entity.NewMail(*recipient, input.Title, buff.Bytes())

	if err := e.emailProvider.Send(ctx, mail); err != nil {
		return err
	}

	return nil
}
