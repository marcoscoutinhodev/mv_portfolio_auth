package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/marcoscoutinhodev/ms_email_notification/config"
	"github.com/marcoscoutinhodev/ms_email_notification/internal/infra"
	"github.com/marcoscoutinhodev/ms_email_notification/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

func init() {
	config.Load()
}

func main() {
	emailService := service.NewEmail(infra.NewEmailProvider(
		config.GMAIL_IDENTITY,
		config.GMAIL_HOST,
		config.GMAIL_PORT,
		config.GMAIL_USER,
		config.GMAIL_SECRET,
	))

	consumer := infra.NewConsumer()
	defer consumer.Close()

	delivery, err := consumer.Delivery()
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	for msg := range delivery {
		go func(m amqp.Delivery) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := emailService.AuthNotification(ctx, m.Body); err != nil {
				fmt.Printf("error: %v", err)
				return
			}

			m.Ack(false)
		}(msg)
	}

	wg.Wait()
}
