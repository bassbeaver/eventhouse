# Image-stage to build application for PROD environment or to run in DEV environment.
FROM golang:1.16-alpine3.13 AS compiler

ARG SERVICE_ENV="dev"
ARG GO_MOD_CACHE_PATH

ENV CGO_ENABLED=0 \
    SERVICE_ENV=${SERVICE_ENV}

RUN apk update && apk upgrade && \
    apk add git bash && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    apk add protobuf && \
    go get -u github.com/golang/protobuf/protoc-gen-go

COPY src /app/src
COPY load-testing /app/load-testing
COPY docker/images/app /app/src/docker/images/app
COPY docker/volumes/gopath /app/src/docker/volumes/gopath

    # Creating entrypoint folder
RUN mkdir -p /scripts && \
    cp -R /app/src/docker/images/app/scripts/* /scripts && \
    chmod +x /scripts/*.sh && \
    # Running build script
    /bin/bash /scripts/build.sh

ENTRYPOINT ["/scripts/entrypoint.sh"]


# Image to run in PROD environment. Thin and without unneeded things like compiler, debugger etc.
FROM alpine:3.13 AS dist

RUN apk update && apk upgrade && \
    apk add bash

# Copying built binaries and other things
COPY --from=compiler ["/app/dist/",  "/app/dist/"]
# Copying entrypoint
COPY --from=compiler ["/scripts/entrypoint.sh", "/scripts/"]

ENTRYPOINT ["/scripts/entrypoint.sh"]