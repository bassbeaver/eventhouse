#!/usr/bin/env bash

ENV=$1
BUILD_IMAGE=$2

docker-compose stop app

if [[ $ENV = "prod" ]]
then
  echo 'Running PROD env'
  sed -i -e 's/^SERVICE_ENV=.*$/SERVICE_ENV=prod/g' app.env

  if [[ -n $BUILD_IMAGE ]]
  then
    echo 'Building image'
    docker build --build-arg SERVICE_ENV=prod --build-arg GO_MOD_CACHE_PATH=/app/src/docker/volumes/gopath/pkg/mod --file images/app/Dockerfile --force-rm --tag bassbeaver/eventhouse:latest ../. # paths are specified relative to the Makefile
  fi
else
  echo 'Running DEV env'
  sed -i -e 's/^SERVICE_ENV=.*$/SERVICE_ENV=dev/g' app.env

  if [[ -n $BUILD_IMAGE ]]
  then
    echo 'Building image'
    docker build --build-arg SERVICE_ENV=dev --file images/app/Dockerfile --force-rm --tag bassbeaver/eventhouse:latest --target compiler ../. # paths are specified relative to the Makefile
  fi
fi

docker-compose up -d
