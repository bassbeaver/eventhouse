FROM golang:1.16-alpine3.13

RUN apk update && apk upgrade && \
    apk add git bash && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    apk add protobuf && \
    go get -u github.com/golang/protobuf/protoc-gen-go
