package grpc

import (
	"errors"
	"fmt"
	apiEvent "github.com/bassbeaver/eventhouse/api/compiled/event"
	"github.com/bassbeaver/eventhouse/service/auth"
	"github.com/bassbeaver/eventhouse/service/logger"
	"github.com/bassbeaver/eventhouse/service/opentracing"
	"github.com/bassbeaver/eventhouse/service/recovery"
	"github.com/bassbeaver/eventhouse/service/request_id_setter"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcMiddlewareAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpcMiddlewareRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

const (
	GrpcServerServiceAlias = "GrpcServer"
)

type GrpcServer struct {
	*grpc.Server
	port int
}

func (s *GrpcServer) Serve() {
	listenAddress := fmt.Sprintf(":%d", s.port)
	tcpListener, tcpListenerError := net.Listen("tcp", listenAddress)
	if tcpListenerError != nil {
		panic(errors.New(fmt.Sprintf("error listening tcp on port: %+v", tcpListenerError)))
	}

	reflection.Register(s.Server)
	err := s.Server.Serve(tcpListener)
	fmt.Printf("gRPC listen error: %+v", err)
}

// --------

func NewGrpcServer(
	port int,
	recoveryService *recovery.RecoveryService,
	opentracingBridge *opentracing.Bridge,
	requestIdISetter *request_id_setter.RequestIdSetter,
	requestContextLoggerSetter *logger.RequestContextLoggerSetter,
	authService *auth.AuthService,
	eventSourceServer apiEvent.APIServer,
) *GrpcServer {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				// Panic interceptor used to recover panic occurred during work of primary interceptors like requestIdISetter and requestContextLoggerSetter
				grpcMiddlewareRecovery.UnaryServerInterceptor(
					grpcMiddlewareRecovery.WithRecoveryHandlerContext(recoveryService.Recover),
				),
				opentracingBridge.Intercept,
				requestIdISetter.Intercept,
				requestContextLoggerSetter.Intercept,
				// Main panic interceptor. Set up after requestIdISetter and requestContextLoggerSetter to make that services available during recovery process
				grpcMiddlewareRecovery.UnaryServerInterceptor(
					grpcMiddlewareRecovery.WithRecoveryHandlerContext(recoveryService.Recover),
				),
				grpcMiddlewareAuth.UnaryServerInterceptor(authService.Auth),
			),
		),
		grpc.StreamInterceptor(
			grpcMiddleware.ChainStreamServer(
				// Panic interceptor used to recover panic occurred during work of primary interceptors like requestIdISetter and requestContextLoggerSetter
				grpcMiddlewareRecovery.StreamServerInterceptor(
					grpcMiddlewareRecovery.WithRecoveryHandlerContext(recoveryService.Recover),
				),
				opentracingBridge.InterceptStream,
				requestIdISetter.InterceptStream,
				requestContextLoggerSetter.InterceptStream,
				// Main panic interceptor. Set up after requestIdISetter and requestContextLoggerSetter to make that services available during recovery process
				grpcMiddlewareRecovery.StreamServerInterceptor(
					grpcMiddlewareRecovery.WithRecoveryHandlerContext(recoveryService.Recover),
				),
				grpcMiddlewareAuth.StreamServerInterceptor(authService.Auth),
			),
		),
	)

	apiEvent.RegisterAPIServer(grpcServer, eventSourceServer)

	return &GrpcServer{port: port, Server: grpcServer}
}
