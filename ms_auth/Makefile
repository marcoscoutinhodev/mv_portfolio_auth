ENV := $(PWD)/.env
include $(ENV)

migrate:
	migrate -path=./migration -database=$(DB_URI) -verbose up

migratedown:
	migrate -path=./migration -database=$(DB_URI) -verbose down

build:
	docker-compose -f docker-compose.prd.yaml build

.PHONY: migrate migratedown build
