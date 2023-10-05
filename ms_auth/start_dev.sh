#!/bin/bash

docker compose -f docker-compose.dev.yaml up -d && \
  docker exec -it ms_auth_app bash
