package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"text/template"

	"github.com/marcoscoutinhodev/ms_auth/config"
	"github.com/marcoscoutinhodev/ms_auth/internal/entity"
	"github.com/marcoscoutinhodev/ms_auth/internal/infra/amqp"
)

type EmailNotification struct {
	amqp                      amqp.AMQP
	EmailNotificationKey      string
	EmailNotificationExchange string
}

func NewEmailNotification() *EmailNotification {
	return &EmailNotification{
		amqp:                      *amqp.NewAMQP(),
		EmailNotificationKey:      config.EMAIL_NOTIFICATION_KEY,
		EmailNotificationExchange: config.EMAIL_NOTIFICATION_EXCHANGE,
	}
}

type emailNotificationPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Title    string `json:"title"`
	Template string `json:"template"`
}

func (q EmailNotification) send(ctx context.Context, user *entity.User, urlBase, token, layout, title string) error {
	payload := struct {
		entity.User
		URLPath string
	}{
		User:    *user,
		URLPath: urlBase + token,
	}

	var buff bytes.Buffer

	t := template.Must(template.New("Template").Parse(layout))
	if err := t.Execute(&buff, payload); err != nil {
		return err
	}

	enp := emailNotificationPayload{
		Name:     user.Name,
		Email:    user.Email,
		Title:    title,
		Template: buff.String(),
	}

	body, err := json.Marshal(&enp)
	if err != nil {
		return err
	}

	if err := q.amqp.Producer(ctx, q.EmailNotificationKey, q.EmailNotificationExchange, body); err != nil {
		return err
	}

	return nil
}

func (q EmailNotification) Register(ctx context.Context, user *entity.User, token string) error {
	layout := `	<p>Hey {{ .Name }}, welcome to MV's Portfolio</p>
							<p>To confirm your email address <a href="{{ .URLPath }}">click here</a></p>`
	title := "Confirm your email address"

	if err := q.send(ctx, user, config.REGISTER_NOTIFICATION_URL, token, layout, title); err != nil {
		return err
	}

	return nil
}

func (q EmailNotification) ForgottenPassword(ctx context.Context, user *entity.User, token string) error {
	layout := `	<p>Hey {{ .Name }}, We've received your forgotten password request.</p>
							<p>To change your password <a href="{{ .URLPath }}">click here</a></p>
							<p>If you didn't request it, don't worry, you can just ignore it, as long as you never share your email and password, your account will be safe</p>`
	title := "Recover your password"

	if err := q.send(ctx, user, config.FORGOT_PASSWORD_NOTIFICATION_URL, token, layout, title); err != nil {
		return err
	}

	return nil
}
