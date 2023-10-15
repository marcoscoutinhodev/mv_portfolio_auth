package service

import (
	"context"

	"github.com/marcoscoutinhodev/ms_email_notification/internal/entity"
)

type EmailProviderInterface interface {
	Send(ctx context.Context, mail *entity.Mail) error
}
