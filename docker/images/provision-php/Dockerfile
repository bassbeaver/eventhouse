FROM golang:1.16

ARG PROTOBUF_VERSION=3.17.3
ARG PHP_GRPC_VERSION=1.39.0

# Utils
RUN apt-get update && \
    apt-get install unzip && \
# Protobuf
    mkdir -p /protobuf && \
    cd /protobuf && \
    wget https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip -O protobuf.zip && \
    unzip protobuf.zip && rm protobuf.zip && \
# grpc PHP (generate client)
    apt-get install php php-dev php-pear phpunit zlib1g-dev -y && \
    pecl install grpc-${PHP_GRPC_VERSION} && \
    cd /protobuf && git clone -b v${PHP_GRPC_VERSION} https://github.com/grpc/grpc && \
    cd /protobuf/grpc && git submodule update --init && \
    apt-get install autoconf libtool automake build-essential -y

# Building protoc & plugins in another layer
RUN apt-get install cmake -y && \
    cd /protobuf/grpc && \
    mkdir -p cmake/build && \
    cd cmake/build && \
    cmake ../.. && \
    make protoc grpc_php_plugin