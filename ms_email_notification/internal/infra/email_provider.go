package infra

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/marcoscoutinhodev/ms_email_notification/internal/entity"
)

type EmailProvider struct {
	Identity string
	Host     string
	Port     string
	User     string
	Password string
	auth     smtp.Auth
	Next     *EmailProvider
}

func NewEmailProvider(identity, host, port, user, password string) *EmailProvider {
	auth := smtp.PlainAuth(identity, user, password, host)
	return &EmailProvider{
		Identity: identity,
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		auth:     auth,
	}
}

func (e *EmailProvider) AddProvider(emailProvider *EmailProvider) *EmailProvider {
	e.Next = emailProvider
	return emailProvider
}

func (e EmailProvider) Send(ctx context.Context, mail *entity.Mail) error {
	to := []string{mail.Recipient.Email}
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"+
		"%s\r\n", mail.Recipient.Email, mail.Title, string(mail.Body)))

	if err := smtp.SendMail(fmt.Sprintf("%s:%s", e.Host, e.Port), e.auth, e.User, to, msg); err != nil {
		if e.Next != nil {
			return e.Next.Send(ctx, mail)
		}

		return err
	}

	return nil
}
