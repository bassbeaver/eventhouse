#!/usr/bin/env bash

SERVICE=$1

if [[ -z $SERVICE ]]
then
  echo "Service to debug not specified"
  exit 1
fi

echo "Starting debug of ${SERVICE}"

docker-compose stop app projector-go

if [[ $SERVICE = "app" ]]
then
  sed -i -e 's/^SERVICE_ENV=.*$/SERVICE_ENV=debug/g' app.env
  sed -i -e 's/^SERVICE_ENV=.*$/SERVICE_ENV=dev/g' projector-go.env
elif [[ $SERVICE = "projector-go" ]]
then
  sed -i -e 's/^SERVICE_ENV=.*$/SERVICE_ENV=dev/g' app.env
  sed -i -e 's/^SERVICE_ENV=.*$/SERVICE_ENV=debug/g' projector-go.env
else
  echo "Unknown service ${SERVICE}"
  exit 1
fi

docker-compose up -d