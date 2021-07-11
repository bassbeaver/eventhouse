#!/usr/bin/env bash

ENV=$1

if [[ $ENV = "prod" ]]
then
  echo 'Building image for PROD env'
  docker build --build-arg SERVICE_ENV=prod --build-arg GO_MOD_CACHE_PATH=/app/src/docker/volumes/gopath/pkg/mod --file images/app/Dockerfile --force-rm --tag bassbeaver/eventhouse:latest ../. # paths are specified relative to the Makefile
else
  echo 'Building image for DEV env'
  docker build --build-arg SERVICE_ENV=dev --file images/app/Dockerfile --force-rm --tag bassbeaver/eventhouse:latest --target compiler ../. # paths are specified relative to the Makefile
fi
