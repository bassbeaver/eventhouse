#!/bin/bash

set -e

if [[ "prod" != $SERVICE_ENV ]]
then
    echo "go version: $(go version)"
    echo "protoc version: $(protoc --version)"

    echo "Starting projector compilation"
    /bin/bash /scripts/compile.sh
fi

echo "Starting projector at $SERVICE_ENV environment"

cd /app/dist # this is required for process to get correct $PWD

if [[ "prod" = $SERVICE_ENV ]]
then
    exec /app/dist/projector -config=/app/dist/config
elif [[ "debug" = $SERVICE_ENV ]]
then
    exec /go/bin/dlv --listen=:40000 --headless=true --log --api-version=2 --accept-multiclient exec /app/dist/projector -- -config=/app/dist/config
else
    exec /app/dist/projector -config=/app/dist/config
fi