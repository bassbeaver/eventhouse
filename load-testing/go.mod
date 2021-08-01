module github.com/bassbeaver/eventhouse/load-testing

go 1.14

require (
	github.com/bassbeaver/eventhouse v0.0.0-00010101000000-000000000000
	github.com/bojand/ghz v0.96.0
	github.com/golang/protobuf v1.5.2
	github.com/jhump/protoreflect v1.5.0
	google.golang.org/grpc v1.39.0
)

replace github.com/bassbeaver/eventhouse => ../src
