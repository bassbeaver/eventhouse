# Based on https://github.com/grpc/grpc/blob/master/examples/php/echo/cli.Dockerfile

FROM composer:1.8.6 as composer

FROM php:7.4


ARG PHP_GRPC_VERSION=1.39.0

RUN apt-get -qq update && \
    apt-get -qq install -y git zlib1g-dev && \
    pecl install grpc-${PHP_GRPC_VERSION} && docker-php-ext-enable grpc

WORKDIR /app

COPY --from=composer /usr/bin/composer /usr/bin/composer

COPY docker/images/php/protoc /usr/local/bin/protoc
COPY docker/images/php/grpc_php_plugin /usr/local/bin/grpc_php_plugin
COPY src/projector/php /app

CMD ["/bin/bash"]