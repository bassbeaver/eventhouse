# EventHouse

Test bench project to test event store application concept for event based systems.

Event store consists of golang application providing gRPC API and column-oriented DBMS (Yandex Clickhouse) to store events. 

## Development

To run project in development mode (all command have to be run from `./docker` directory):

1. Create `app.env` file:
   ```shell
   cp app.env.example app.env
   ```
2. Build basic docker image:
    ```shell
    make build
    ```
3. Create `./src/config/config-secret.yaml`:
   ```shell
   cp ../src/config/config-secret.yaml.example ../src/config/config-secret.yaml
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
docker-compose run --rm provision protoc --proto_path=/app/src/api/proto --go_out=plugins=grpc:/app/src/api/compiled <filename>.proto
```


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