#!/bin/bash

set -e

if [[ "prod" != $SERVICE_ENV ]]
then
    go version

    echo "Starting application compilation"
    /bin/bash /scripts/compile.sh
fi

echo "Starting application at $SERVICE_ENV environment"

cd /app/dist # this is required for process to get correct $PWD

if [[ "prod" = $SERVICE_ENV ]]
then
    exec /app/dist/eventhouse -config=/app/dist/config
elif [[ "debug" = $SERVICE_ENV ]]
then
    exec /go/bin/dlv --listen=:40000 --headless=true --log --api-version=2 --accept-multiclient exec /app/dist/eventhouse -- -config=/app/dist/config
else
    exec /app/dist/eventhouse -config=/app/dist/config
fi