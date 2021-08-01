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

echo "Application compilation for $SERVICE_ENV environment"

cd /app/src

if [[ "prod" = $SERVICE_ENV ]]
then
    go build -o /app/dist/eventhouse
else
    go build -gcflags "all=-N -l" -o /app/dist/eventhouse
fi

echo "Application compiled"

chmod 555 /app/dist/eventhouse

rm -rf /app/dist/config

cp -R /app/src/config /app/dist
