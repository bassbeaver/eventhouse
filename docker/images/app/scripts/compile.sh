#!/usr/bin/env bash

set -e

if [[ -d ${GOPATH}/pkg/mod ]]
then
    echo "Golang modules cache found: ${GOPATH}/pkg/mod"
elif [[ -d $GO_MOD_CACHE_PATH ]]
then
    echo "Golang cache path parameter is set and cache is available at ${GO_MOD_CACHE_PATH}. Creating symlink"

    if [[ ! -d ${GOPATH}/pkg ]]
    then
      mkdir ${GOPATH}/pkg
    fi

    ln -s $GO_MOD_CACHE_PATH ${GOPATH}/pkg/mod
fi


mkdir -p /app/dist

echo "Compiling load testing tools"

cd /app/load-testing
go build -o /app/dist/load-testing

echo "Application compilation for $SERVICE_ENV environment"

if [[ "prod" = $SERVICE_ENV ]]
then
    cd /app/src
    go build -o /app/dist/eventhouse
else
    cd /app/src
    go build -gcflags "all=-N -l" -o /app/dist/eventhouse
fi

echo "Application compiled"

chmod 555 /app/dist/eventhouse

rm -rf /app/dist/config

cp -R /app/src/config /app/dist
