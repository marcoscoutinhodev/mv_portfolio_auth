package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	DB_DRIVER                             string
	DB_URI                                string
	DB_NAME                               string
	DB_MAX_IDLE_CONNS                     string
	SERVER_PORT                           string
	RBMQ_URI                              string
	REGISTER_NOTIFICATION_KEY             string
	REGISTER_NOTIFICATION_EXCHANGE        string
	FORGOT_PASSWORD_NOTIFICATION_KEY      string
	FORGOT_PASSWORD_NOTIFICATION_EXCHANGE string
	JWT_PRIVATE_KEY                       string
	JWT_PUBLIC_KEY                        string
	JWT_SECRET_KEY                        string
)

func Load() {
	// there is no need to check the error because in PRD the application does not use the .env file
	godotenv.Load()

	var errors []string

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

	if REGISTER_NOTIFICATION_KEY = os.Getenv("REGISTER_NOTIFICATION_KEY"); REGISTER_NOTIFICATION_KEY == "" {
		errors = append(errors, "REGISTER_NOTIFICATION_KEY not found")
	}

	if REGISTER_NOTIFICATION_EXCHANGE = os.Getenv("REGISTER_NOTIFICATION_EXCHANGE"); REGISTER_NOTIFICATION_EXCHANGE == "" {
		errors = append(errors, "REGISTER_NOTIFICATION_EXCHANGE not found")
	}

	if FORGOT_PASSWORD_NOTIFICATION_KEY = os.Getenv("FORGOT_PASSWORD_NOTIFICATION_KEY"); FORGOT_PASSWORD_NOTIFICATION_KEY == "" {
		errors = append(errors, "FORGOT_PASSWORD_NOTIFICATION_KEY not found")
	}

	if FORGOT_PASSWORD_NOTIFICATION_EXCHANGE = os.Getenv("FORGOT_PASSWORD_NOTIFICATION_EXCHANGE"); FORGOT_PASSWORD_NOTIFICATION_EXCHANGE == "" {
		errors = append(errors, "FORGOT_PASSWORD_NOTIFICATION_EXCHANGE not found")
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
