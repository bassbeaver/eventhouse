module github.com/bassbeaver/eventhouse/projector

go 1.16

require (
	github.com/bassbeaver/eventhouse v0.0.0-00010101000000-000000000000
	github.com/bassbeaver/gioc v0.1.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/spf13/viper v1.8.1
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	google.golang.org/grpc v1.38.0
)

replace github.com/bassbeaver/eventhouse => ../../src
