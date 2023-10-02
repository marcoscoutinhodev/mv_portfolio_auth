#!/bin/bash

docker compose -f docker-compose.dev.yaml up -d && \
  docker exec -it mv_chat_app bash
