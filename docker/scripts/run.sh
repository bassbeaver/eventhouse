#!/usr/bin/env bash

ENV=$1

docker-compose stop app

if [[ $ENV = "prod" ]]
then
  echo 'Running PROD env'
  sed -i -e 's/^SERVICE_ENV=.*$/SERVICE_ENV=prod/g' app.env projector-go.env

else
  echo 'Running DEV env'
  sed -i -e 's/^SERVICE_ENV=.*$/SERVICE_ENV=dev/g' app.env projector-go.env
fi

docker-compose up -d
