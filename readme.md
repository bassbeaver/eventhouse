# EventHouse

Test bench project to test event store application concept for event based systems.

Event store consists of golang application providing gRPC API and column-oriented DBMS (Yandex Clickhouse) to store events. 

## Development

To run project in development mode (all command have to be run from `./docker` directory):

1. Create `app.env` file for main event store application:
   ```shell
   cp app.env.example app.env
   ```

   Create `projector-go.env` file for golang version of events projector:
   ```shell
   cp projector-go.env.example projector-go.env
   ```
2. Build basic docker image:
    ```shell
    make build
    ```
3. Create local part of config for event store`./src/config/config-secret.yaml`:
   ```shell
   cp ../src/config/config-secret.yaml.example ../src/config/config-secret.yaml
   ```
   Create local part of config for golang version of event projector `./projector/go/config/config-secret.yaml`:
   ```shell
   cp ../projector/go/config/config-secret.yaml.example ../src/config/config-secret.yaml
   ```
4. Start project:
    ```shell
    make run dev
    ```
5. Create API clients (see [Create API clients](#create-api-clients) section)
6. Run API requests :)

#### Debug mode

To run project in debug mode, instead of `make run dev` use `make debug app` command, 
and do not forget to set up debugger listener in your IDE. Debugger listener have to connect to `192.168.66.10:40000`.

#### Adding go mod dependencies

From `./docker` directory:

```shell script
docker-compose run --rm --workdir="/app/src" provision go get github.com/some/go/package
```

Just download packages listed in go.mod/go.sum files (without build / install runs):
```shell script
docker-compose run --rm --workdir="/app/src" provision go get -d
```

#### Compiling Protobuf

To compile protobuf files for Golang run:
```shell script
docker-compose run --rm provision protoc --proto_path=/app/src/api/proto --go_out=plugins=grpc:/app/src/api/compiled event.proto
```

To compile protobuf files for PHP run:
```shell script
docker-compose run --rm php protoc \
    --proto_path=/api-proto \
    --php_out=/app/api \
    --grpc_out=/app/api \
    --plugin=protoc-gen-grpc=/usr/local/bin/grpc_php_plugin \
    event.proto
```

#### Compiling protoc and gRPC PHP plugin

For some reasons there are no already compiled gRPC PHP plugin on the internet so you have to compile it yourself.
I've compiled it for you located it in `docker/images/php/grpc_php_plugin`. 
Also, I've compiled protoc compiler (located in `docker/images/php/protoc`).

But if for some reasons you want to recompile it again - uncomment `provision-php` service in docker-compose.yml and build
it. Compiled binaries will be available at:
* protoc: `/protobuf/grpc/cmake/build/third_party/protobuf/protoc`
* gRPC PHP plugin: `/protobuf/grpc/cmake/build/grpc_php_plugin`

### Production

1. Build basic docker image (run from `./docker` directory):
    ```shell
    make build prod
    ```
2. Create `config-secret.yaml` file. You can use `./src/config/config-secret.yaml.example` as an example, 
   but do not forget to change config values to appropriate :) Copy created file to you server.
3. Deploy `bassbeaver/eventhouse:latest` docker image to your server.
4. Run application container. 
   
   Do not forget to link `config-secret.yaml` file into application container to `/app/dist/config` folder (with docker volume).

   Run command can look like:
   ```shell
   docker run -d --name eventhouse --volume=/path/to/config-secret.yaml:/app/dist/config/config-secret.yaml
   ```

### Load testing

Project already have some code to perform load testing of event store. Load testing code is located in `./load-testing` folder
and uses [https://ghz.sh](https://ghz.sh) golang package to perform load. Load testing binary compiles together with main application
and is located in `/app/dist/load-testing` inside container.

### Create API clients

To create two clients accounts with credentials:
```shell
client1 secret1
client2 secret2
```

Connect to Clickhouse and run next query:
```shell
 INSERT INTO eventhouse.apiClients(ClientId, SecretHash) 
 VALUES 
  ('client1', '$2a$04$YOO1B9PjJ2iYv6ygasD7WOmSeTZ14oexV8gw9RBwMIC2.kzsvJ3iu'),
  ('client2', '$2a$04$EOPyuUupZDfDxdaZXwW8pOLcL.Xxq.TviWpgGUTUfxH2FLh1LWdTC')
```

### Sending example events

To send gRPC requests from console we can use gRPCurl: https://github.com/fullstorydev/grpcurl

```bash
grpcurl --plaintext \
  -H 'Authorization: Basic Y2xpZW50MTpzZWNyZXQx' \
  -d '{"idempotencyKey": "ik1", "eventType": "Created", "entityType": "Subscription", "entityId": "sub_1", "payload": "{\"amount\":15.5,\"plan\":{\"level\":\"basic\",\"duration\":\"1m\"},\"transaction\":{\"id\":\"tr_1\",\"amount\":15.5}}"}' \
  192.168.66.10:750 eventhouse.grpc.event.API/Push
```

```bash
grpcurl --plaintext \
  -H 'Authorization: Basic Y2xpZW50MTpzZWNyZXQx' \
  -d '{"idempotencyKey": "ik2", "eventType": "Renewed", "entityType": "Subscription", "entityId": "sub_1", "payload": "{\"transaction\":{\"id\":\"tr_2\",\"amount\":15.5}}"}' \
  192.168.66.10:750 eventhouse.grpc.event.API/Push
```

```bash
grpcurl --plaintext \
  -H 'Authorization: Basic Y2xpZW50MTpzZWNyZXQx' \
  -d '{"idempotencyKey": "ik3", "eventType": "Canceled", "entityType": "Subscription", "entityId": "sub_1"}' \
  192.168.66.10:750 eventhouse.grpc.event.API/Push
```
