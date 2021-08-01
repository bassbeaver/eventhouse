# Image-stage to build application for PROD environment or to run in DEV environment.
FROM golang:1.16-alpine3.13 AS compiler

ARG SERVICE_ENV="dev"

ENV CGO_ENABLED=0 \
    SERVICE_ENV=${SERVICE_ENV} \
    VOLUME_GOPATH="/app/src/docker/volumes/gopath"

RUN apk update && apk upgrade && \
    apk add git bash && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    apk add protobuf && \
    go get -u github.com/golang/protobuf/protoc-gen-go && \
    # Removing $GOPATH cache to be able to use cache from $VOLUME_GOPATH (we do not know how to overlay them, so just remove cache created from delve and protoc installation) \
    # Also, installation of delve and protoc have to be in basic image in case of good enterprise architecture
    rm -R ${GOPATH}/pkg/mod && \
    echo "go version: $(go version)" && \
    echo "protoc version: $(protoc --version)"

COPY src /app/src
COPY docker/images/app /app/src/docker/images/app
COPY docker/volumes/gopath $VOLUME_GOPATH

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