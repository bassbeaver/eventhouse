package opentracing

import (
	"context"
	"fmt"
	loggerService "github.com/bassbeaver/eventhouse/service/logger"
	"github.com/bassbeaver/logopher"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	grpcMetadata "google.golang.org/grpc/metadata"
	"io"
)

const (
	OpentracingBridgeServiceAlias = "OpentracingBridgeService"
	methodSubscribeGlobalStream   = "/eventhouse.grpc.event.API/SubscribeGlobalStream"
)

type Bridge struct {
	loggerObj    *logopher.Logger
	tracer       opentracing.Tracer
	tracerCloser io.Closer
}

func (br *Bridge) Intercept(
	contextObj context.Context,
	requestObj interface{},
	serverInfo *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	rootSpan := br.tracer.StartSpan(
		serverInfo.FullMethod,
		ext.RPCServerOption(br.extractExternalOpentracingContext(contextObj, serverInfo.FullMethod)),
	)
	defer rootSpan.Finish()

	responseObj, responseError := handler(opentracing.ContextWithSpan(contextObj, rootSpan), requestObj)

	return responseObj, responseError
}

func (br *Bridge) InterceptStream(
	serviceObj interface{},
	serverStreamObj grpc.ServerStream,
	serverInfo *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	rootSpan := br.tracer.StartSpan(
		serverInfo.FullMethod,
		ext.RPCServerOption(br.extractExternalOpentracingContext(serverStreamObj.Context(), serverInfo.FullMethod)),
	)

	// For long living "subscriptions" requests we recognize root span as "request beginning" mark and finish it immediately.
	// Also, there will be another "closing" span marking end of request.
	if methodSubscribeGlobalStream == serverInfo.FullMethod {
		rootSpan.Finish()
		defer func() {
			endingSpan := br.tracer.StartSpan(
				serverInfo.FullMethod+":finish",
				ext.RPCServerOption(br.extractExternalOpentracingContext(serverStreamObj.Context(), serverInfo.FullMethod)),
			)
			endingSpan.Finish()
		}()
	} else {
		defer rootSpan.Finish()
	}

	newCtx := opentracing.ContextWithSpan(serverStreamObj.Context(), rootSpan)

	wrappedServerStreamObj := grpc_middleware.WrapServerStream(serverStreamObj)
	wrappedServerStreamObj.WrappedContext = newCtx

	return handler(serviceObj, wrappedServerStreamObj)
}

// CloseTracer closes tracer object owned by Bridge
// TODO implement graceful shutdown and use Bridge.CloseTracer() inside it
func (br *Bridge) CloseTracer() {
	if nil != br.tracerCloser {
		closeError := br.tracerCloser.Close()
		if nil != closeError {
			br.loggerObj.Critical("Failed to close JaeggerTracer", &logopher.MessageContext{"error": closeError.Error()})
		}
	}
}

func (br *Bridge) extractExternalOpentracingContext(contextObj context.Context, grpcMethodName string) opentracing.SpanContext {
	metadata, metadataIsOk := grpcMetadata.FromIncomingContext(contextObj)
	if !metadataIsOk {
		br.loggerObj.Error(
			"Failed to extract metadata from gRPC context",
			&logopher.MessageContext{"method": grpcMethodName},
		)
	}

	externalOpentracingContext, externalOpentracingContextErr := br.tracer.Extract(
		GrpcMetadata,
		GrpcMetadataCarrier(metadata),
	)
	if nil != externalOpentracingContextErr && opentracing.ErrSpanContextNotFound != externalOpentracingContextErr {
		br.loggerObj.Error(
			"Failed to extract Opentracing context from gRPC metadata",
			&logopher.MessageContext{"method": grpcMethodName, "error": externalOpentracingContextErr.Error()},
		)
	}

	return externalOpentracingContext
}

func (br *Bridge) Tracer() opentracing.Tracer {
	return br.tracer
}

func NewBridge(agentHost string, agentPort int, serviceName string, loggerFactory *loggerService.LoggerFactory) *Bridge {
	loggerObj, _ := loggerFactory.CreateLogger(1)

	jaeggerConfigObj := jaegerConfig.Configuration{
		ServiceName: serviceName,
		Reporter: &jaegerConfig.ReporterConfig{
			LocalAgentHostPort: fmt.Sprintf("%s:%d", agentHost, agentPort),
		},
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
	}

	tracer, tracerCloser, tracerError := jaeggerConfigObj.NewTracer()
	if nil != tracerError {
		loggerObj.Critical("Failed to create JaeggerTracer", &logopher.MessageContext{"error": tracerError.Error()})

		tracer = &opentracing.NoopTracer{}
	}

	return &Bridge{
		loggerObj:    loggerObj,
		tracer:       tracer,
		tracerCloser: tracerCloser,
	}
}
