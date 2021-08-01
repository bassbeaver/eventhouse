#!/bin/bash

set -e

echo "Building image for $SERVICE_ENV environnent"

if [[ "prod" = $SERVICE_ENV ]]
then
    echo "Starting projector compilation"
    /bin/bash /scripts/compile.sh
else
    echo "Projector will be compiled during container start"
fi