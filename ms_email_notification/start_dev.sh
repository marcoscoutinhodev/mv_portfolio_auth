#!/bin/bash

docker compose -f docker-compose.dev.yaml up -d && \
  docker exec -it ms_email_notification_app bash
