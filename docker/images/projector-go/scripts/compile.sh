#!/usr/bin/env bash

set -e

if [[ -d ${GOPATH}/pkg/mod ]]
then
    echo "Golang modules cache found: ${GOPATH}/pkg/mod"
elif [[ -d $VOLUME_GOPATH/pkg/mod ]]
then
    echo "Golang cache path parameter is set and cache is available at ${VOLUME_GOPATH}/pkg/mod. Creating symlink"

    if [[ ! -d ${GOPATH}/pkg ]]
    then
      mkdir ${GOPATH}/pkg
    fi

    ln -s ${VOLUME_GOPATH}/pkg/mod ${GOPATH}/pkg/mod
fi


mkdir -p /app/dist

echo "Compiling load testing tools"

cd /app/load-testing
go build -o /app/dist/load-testing

echo "Compiling projector for $SERVICE_ENV environment"

cd /app/projector/go

if [[ "prod" = $SERVICE_ENV ]]
then
    go build -o /app/dist/projector
else
    go build -gcflags "all=-N -l" -o /app/dist/projector
fi

echo "Projector compiled"

chmod 555 /app/dist/projector

rm -rf /app/dist/config

cp -R /app/projector/go/config /app/dist
