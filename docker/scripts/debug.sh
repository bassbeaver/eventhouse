#!/usr/bin/env bash

SERVICE=$1
BUILD_IMAGE=$2

if [[ -z $SERVICE ]]
then
  echo "Service to debug not specified"
  exit 1
fi

echo "Starting debug of ${SERVICE}"

docker-compose stop app

if [[ $SERVICE = "app" ]]
then
  sed -i -e 's/^SERVICE_ENV=.*$/SERVICE_ENV=debug/g' app.env

  if [[ -n $BUILD_IMAGE ]]
  then
    echo 'Building image'
    docker build --build-arg SERVICE_ENV=debug --file images/app/Dockerfile --force-rm --tag bassbeaver/eventhouse:latest --target compiler ../. # paths are specified relative to the Makefile
  fi
else
  echo "Unknown service ${SERVICE}"
  exit 1
fi

docker-compose up -d