package config

import (
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	REGISTER_NOTIFICATION_URL        string
	FORGOT_PASSWORD_NOTIFICATION_URL string
	DB_DRIVER                        string
	DB_URI                           string
	DB_NAME                          string
	DB_MAX_IDLE_CONNS                string
	SERVER_PORT                      string
	RBMQ_URI                         string
	EMAIL_NOTIFICATION_KEY           string
	EMAIL_NOTIFICATION_EXCHANGE      string
	JWT_PRIVATE_KEY                  string
	JWT_PUBLIC_KEY                   string
	JWT_SECRET_KEY                   string
)

func Load() {
	// there is no need to check the error because in PRD the application does not use the .env file
	godotenv.Load()

	var errors []string

	if REGISTER_NOTIFICATION_URL = os.Getenv("REGISTER_NOTIFICATION_URL"); REGISTER_NOTIFICATION_URL == "" {
		errors = append(errors, "REGISTER_NOTIFICATION_URL not found")
	} else {
		if _, err := url.ParseRequestURI(REGISTER_NOTIFICATION_URL + "test"); err != nil {
			errors = append(errors, "REGISTER_NOTIFICATION_URL has invalid format")
		}
	}

	if FORGOT_PASSWORD_NOTIFICATION_URL = os.Getenv("FORGOT_PASSWORD_NOTIFICATION_URL"); FORGOT_PASSWORD_NOTIFICATION_URL == "" {
		errors = append(errors, "FORGOT_PASSWORD_NOTIFICATION_URL not found")
	} else {
		if _, err := url.ParseRequestURI(FORGOT_PASSWORD_NOTIFICATION_URL + "test"); err != nil {
			errors = append(errors, "FORGOT_PASSWORD_NOTIFICATION_URL has invalid format")
		}
	}

	if DB_DRIVER = os.Getenv("DB_DRIVER"); DB_DRIVER == "" {
		errors = append(errors, "DB_DRIVER not found")
	}

	if DB_URI = os.Getenv("DB_URI"); DB_URI == "" {
		errors = append(errors, "DB_URI not found")
	}

	if DB_NAME = os.Getenv("DB_NAME"); DB_NAME == "" {
		errors = append(errors, "DB_NAME not found")
	}

	if DB_MAX_IDLE_CONNS = os.Getenv("DB_MAX_IDLE_CONNS"); DB_MAX_IDLE_CONNS == "" {
		errors = append(errors, "DB_MAX_IDLE_CONNS not found")
	}

	if SERVER_PORT = os.Getenv("SERVER_PORT"); SERVER_PORT == "" {
		errors = append(errors, "SERVER_PORT not found")
	}

	if RBMQ_URI = os.Getenv("RBMQ_URI"); RBMQ_URI == "" {
		errors = append(errors, "RBMQ_URI not found")
	}

	if EMAIL_NOTIFICATION_KEY = os.Getenv("EMAIL_NOTIFICATION_KEY"); EMAIL_NOTIFICATION_KEY == "" {
		errors = append(errors, "EMAIL_NOTIFICATION_KEY not found")
	}

	if EMAIL_NOTIFICATION_EXCHANGE = os.Getenv("EMAIL_NOTIFICATION_EXCHANGE"); EMAIL_NOTIFICATION_EXCHANGE == "" {
		errors = append(errors, "EMAIL_NOTIFICATION_EXCHANGE not found")
	}

	if JWT_PRIVATE_KEY = os.Getenv("JWT_PRIVATE_KEY"); JWT_PRIVATE_KEY == "" {
		errors = append(errors, "JWT_PRIVATE_KEY not found")
	}

	if JWT_PUBLIC_KEY = os.Getenv("JWT_PUBLIC_KEY"); JWT_PUBLIC_KEY == "" {
		errors = append(errors, "JWT_PUBLIC_KEY not found")
	}

	if JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY"); JWT_SECRET_KEY == "" {
		errors = append(errors, "JWT_SECRET_KEY not found")
	}

	if len(errors) > 0 {
		panic("\n" + strings.Join(errors, "\n"))
	}
}
