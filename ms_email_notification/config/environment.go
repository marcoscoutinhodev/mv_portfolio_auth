package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	GMAIL_IDENTITY           string
	GMAIL_HOST               string
	GMAIL_PORT               string
	GMAIL_USER               string
	GMAIL_SECRET             string
	RBMQ_URI                 string
	EMAIL_NOTIFICATION_QUEUE string
)

func Load() {
	// there is no need to check the error because in PRD the application does not use the .env file
	godotenv.Load()

	var errors []string

	GMAIL_IDENTITY = os.Getenv("DB_DRIVER")

	if GMAIL_HOST = os.Getenv("GMAIL_HOST"); GMAIL_HOST == "" {
		errors = append(errors, "GMAIL_HOST not found")
	}

	if GMAIL_PORT = os.Getenv("GMAIL_PORT"); GMAIL_PORT == "" {
		errors = append(errors, "GMAIL_PORT not found")
	}

	if GMAIL_USER = os.Getenv("GMAIL_USER"); GMAIL_USER == "" {
		errors = append(errors, "GMAIL_USER not found")
	}

	if GMAIL_SECRET = os.Getenv("GMAIL_SECRET"); GMAIL_SECRET == "" {
		errors = append(errors, "GMAIL_SECRET not found")
	}

	if RBMQ_URI = os.Getenv("RBMQ_URI"); RBMQ_URI == "" {
		errors = append(errors, "RBMQ_URI not found")
	}

	if EMAIL_NOTIFICATION_QUEUE = os.Getenv("EMAIL_NOTIFICATION_QUEUE"); EMAIL_NOTIFICATION_QUEUE == "" {
		errors = append(errors, "EMAIL_NOTIFICATION_QUEUE not found")
	}

	if len(errors) > 0 {
		panic("\n" + strings.Join(errors, "\n"))
	}
}
